package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"returnone/config"
	"returnone/database/redis"
	"returnone/database/user"
	utils "returnone/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
)

// 24 hours later time(for access token)
var twenty_four_hours_later = time.Now().Add(time.Hour * 24)

// 60 days later time(for refresh tokne)
var sixty_days_later = time.Now().Add(time.Hour * 24 * 60)

func SignUp(c *fiber.Ctx) error {
	var data map[string]string
	// get data from body
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid post request",
			})
	}

	// check if data is valid return error
	request_data_error := utils.RequestDataRequired(data, []string{"email", "password", "user_name"})
	if request_data_error != nil {
		return c.Status(400).JSON(request_data_error)
	}

	if !utils.IsValidUsername(data["user_name"]) {
		return c.Status(400).JSON(utils.RequestValueValid("user name"))
	}

	if !utils.IsValidEmail(data["email"]) {
		return c.Status(400).JSON(utils.RequestValueValid("email address"))
	}

	// check the email has already been used
	if userDatabase.CheckUserEmailExist(data["email"]) != 0 {
		return c.Status(400).JSON(utils.RequestValueInUse("email address"))
	}

	// check the user name has already been used
	if userDatabase.CheckUserNameExist(data["user_name"]) != 0 {
		return c.Status(400).JSON(utils.RequestValueInUse("user name"))
	}

	hash_password, _ := utils.HashPassword(data["password"])

	save_data, save_data_err := userDatabase.CreateUser(data["email"], hash_password, data["user_name"])

	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", save_data_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error creating user", save_data_err))
	}

	//generate Jwt token
	access_token, access_token_err := utils.GenerateJwtToken(save_data.Id, "accessToken", twenty_four_hours_later.Unix())
	refresh_token, refresh_token_err := utils.GenerateJwtToken(save_data.Id, "refreshToken", sixty_days_later.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", refresh_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating refresh token", refresh_token_err))
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", access_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating access token", access_token_err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:    "accessToken",
		Value:   access_token,
		Expires: twenty_four_hours_later,
	}
	refresh_token_cookie := fiber.Cookie{
		Name:    "refreshToken",
		Value:   refresh_token,
		Expires: sixty_days_later,
	}
	c.Cookie(&access_token_cookie)
	c.Cookie(&refresh_token_cookie)

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Sign up successfully",
			"data":    save_data,
		})
}

func LogIn(c *fiber.Ctx) error {
	var data map[string]string
	// get data from body
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	if !utils.IsValidEmail(data["email"]) {
		return c.Status(400).JSON(utils.RequestValueValid("email address"))
	}
	// check if data is valid return error
	request_data_error := utils.RequestDataRequired(data, []string{"email", "password"})
	if request_data_error != nil {
		return c.Status(400).JSON(request_data_error)
	}

	//Get hash password from database and check does email exist
	user_data, user_data_err := userDatabase.GetUserPassword(data["email"])

	if user_data_err == sql.ErrNoRows {
		return c.Status(401).JSON(utils.RequestValueValid("password or email"))
	}
	//check password is correct
	check_password := utils.CheckPasswordHash(data["password"], user_data.Password)

	// if password not correct
	if !check_password {
		return c.Status(401).JSON(utils.RequestValueValid("password or email"))
	}

	//generate Jwt token
	access_token, access_token_err := utils.GenerateJwtToken(user_data.Id, "accessToken", twenty_four_hours_later.Unix())
	refresh_token, refresh_token_err := utils.GenerateJwtToken(user_data.Id, "refreshToken", sixty_days_later.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", refresh_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating refresh token", refresh_token_err))
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", access_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating access token", access_token_err))
	}

	if user_data.Default_2fa == 3 {
		if data["otp"] == "" {
			return c.Status(403).JSON(utils.ErrorMessage("OTP is required", nil))
		}

		valid := totp.Validate(data["otp"], user_data.Totp)

		if !valid {
			return c.Status(403).JSON(utils.ErrorMessage("OTP is not valid", nil))
		}
	}
	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:    "accessToken",
		Value:   access_token,
		Expires: twenty_four_hours_later,
	}
	refresh_token_cookie := fiber.Cookie{
		Name:    "refreshToken",
		Value:   refresh_token,
		Expires: sixty_days_later,
	}
	c.Cookie(&access_token_cookie)
	c.Cookie(&refresh_token_cookie)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Successful login",
	})

}

func GoogleLogin(c *fiber.Ctx) error {
	// generate random base 64 string for the state token
	state_token_key, state_token_key_err := utils.GenerateRandomBase64String()
	state_token, state_token_err := utils.GenerateRandomBase64String()

	// handle errors
	if state_token_key_err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error generate state token", state_token_key_err))
	}
	if state_token_err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error generate state token", state_token_key_err))
	}

	// save state token to redis (for verify vaild)
	save_state_token_err := redis.CreateStringData(state_token_key, state_token, time.Minute*15)

	if save_state_token_err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error saving state token", save_state_token_err))
	}

	// actually i don't really know this code will do what
	url := config.GoogleConfig.GoogleLoginConfig.AuthCodeURL(fmt.Sprintf("%v %v", state_token_key, state_token))

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func GoogleCallBack(c *fiber.Ctx) error {
	// get state from query
	state := c.Query("state")

	// cuz the state is "xxxx xxxx" this first xxxx is key the second is value
	result := strings.Split(state, " ")

	// get state token from redis for verify vaild
	save_token, redis_error := redis.GetStrigData(result[0])
	if redis_error != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error get state token", redis_error))
	}

	// if this function is done than run this (i dont want to run this is the middle because it will take some time)
	defer redis.DeleteStringData(result[0])

	//check the states
	if result[1] != save_token {
		return c.Status(500).JSON(utils.ErrorMessage("States don't match", nil))
	}

	// so i dont know what is this work
	code := c.Query("code")

	googlecon := config.GoogleOauth()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Code-Token Exchange Failed", err))
	}

	// get the user data
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("User data fetch failed", err))
	}

	// byte to map
	user_data_byte, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("JSON parsing failed", err))
	}
	json_str := string(user_data_byte)
	var user_data map[string]interface{}
	err = json.Unmarshal([]byte(json_str), &user_data)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("JSON unmarshal failed", err))
	}

	// get user data from database (verify this user is log in or sign up)
	get_usre_data, get_user_data_error := userDatabase.GetGoogleAccount(user_data["id"].(string))

	//generate Jwt token
	access_token, access_token_err := utils.GenerateJwtToken(user_data["id"].(string), "accessToken", twenty_four_hours_later.Unix())
	refresh_token, refresh_token_err := utils.GenerateJwtToken(user_data["id"].(string), "refreshToken", sixty_days_later.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", refresh_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating refresh token", refresh_token_err))
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", access_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating access token", access_token_err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:    "accessToken",
		Value:   access_token,
		Expires: twenty_four_hours_later,
	}
	refresh_token_cookie := fiber.Cookie{
		Name:    "refreshToken",
		Value:   refresh_token,
		Expires: sixty_days_later,
	}

	// if cant found this user than this user is login
	if get_user_data_error == nil {
		c.Cookie(&access_token_cookie)
		c.Cookie(&refresh_token_cookie)
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Account successfully login",
			"data":    get_usre_data,
		})
	}

	// if this user didn't sign up than create user data
	save_data, save_data_err := userDatabase.CreateUserWithGoogleLogin(user_data["id"].(string), user_data["picture"].(string))
	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", save_data_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error creating user", save_data_err))
	}

	c.Cookie(&access_token_cookie)
	c.Cookie(&refresh_token_cookie)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Account successfully signup",
		"data":    save_data,
	})
}

func GithubLogin(c *fiber.Ctx) error {
	// generate random base 64 string for the state token
	state_token_key, state_token_key_err := utils.GenerateRandomBase64String()
	state_token, state_token_err := utils.GenerateRandomBase64String()

	// handle errors
	if state_token_key_err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error generate state token", state_token_key_err))
	}
	if state_token_err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error generate state token", state_token_key_err))
	}

	// save state token to redis (for verify vaild)
	save_state_token_err := redis.CreateStringData(state_token_key, state_token, time.Minute*15)

	if save_state_token_err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error saving state token", save_state_token_err))
	}

	// actually i don't really know this code will do what
	url := config.GithubConfig.GithubLoginConfig.AuthCodeURL(fmt.Sprintf("%v %v", state_token_key, state_token))

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func GithubCallBack(c *fiber.Ctx) error {
	// get state from query
	state := c.Query("state")

	// cuz the state is "xxxx xxxx" this first xxxx is key the second is value
	result := strings.Split(state, " ")

	// get state token from redis for verify vaild
	save_token, redis_error := redis.GetStrigData(result[0])
	if redis_error != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error get state token", redis_error))
	}

	// if this function is done than run this (i dont want to run this is the middle because it will take some time)
	defer redis.DeleteStringData(result[0])

	//check the states
	if result[1] != save_token {
		return c.Status(500).JSON(utils.ErrorMessage("States don't match", nil))
	}

	// so i dont know what is this work
	code := c.Query("code")

	githubcon := config.GithubOauth()

	token, err := githubcon.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Code-Token Exchange Failed", err))
	}

	//set the request for get user data
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("User data fetch failed", err))
	}

	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	// get the user data
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("User data fetch failed", err))
	}

	// byte to map
	user_data_byte, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("JSON parsing failed", err))
	}
	json_str := string(user_data_byte)
	var user_data map[string]interface{}
	err = json.Unmarshal([]byte(json_str), &user_data)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("JSON unmarshal failed", err))
	}

	// convert user data id to string
	user_github_id := fmt.Sprintf("%f", user_data["id"])

	// get user data from database (verify this user is log in or sign up)
	get_usre_data, get_user_data_error := userDatabase.GetGithubAccount(user_github_id)

	//generate Jwt token
	access_token, access_token_err := utils.GenerateJwtToken(user_github_id, "accessToken", twenty_four_hours_later.Unix())
	refresh_token, refresh_token_err := utils.GenerateJwtToken(user_github_id, "refreshToken", sixty_days_later.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", refresh_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating refresh token", refresh_token_err))
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", access_token_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error generating access token", access_token_err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:    "accessToken",
		Value:   access_token,
		Expires: twenty_four_hours_later,
	}
	refresh_token_cookie := fiber.Cookie{
		Name:    "refreshToken",
		Value:   refresh_token,
		Expires: sixty_days_later,
	}
	// if cant found this user than this user is login
	if get_user_data_error == nil {
		c.Cookie(&access_token_cookie)
		c.Cookie(&refresh_token_cookie)
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Account successfully login",
			"data":    get_usre_data,
		})
	}

	// if this user didn't sign up than create user data
	save_data, save_data_err := userDatabase.CreateUserWithGithubLogin(user_github_id, user_data["avatar_url"].(string))
	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", save_data_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error creating user", save_data_err))
	}

	c.Cookie(&access_token_cookie)
	c.Cookie(&refresh_token_cookie)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Account successfully signup",
		"data":    save_data,
	})
}

func EmailExist(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	if data["email"] == "" {
		return c.Status(400).JSON(utils.RequestValueRequired("Email address"))
	}

	if !utils.IsValidEmail(data["email"]) {
		return c.Status(400).JSON(utils.RequestValueValid("Email address"))
	}

	// Check if the email is already in the database
	if userDatabase.CheckUserEmailExist(data["email"]) != 0 {
		return c.Status(200).JSON(utils.RequestValueInUse("email address"))
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "This email has not been used yet",
			"inuse":   false,
		})
}

func UserNameExist(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	if data["user_name"] == "" {
		return c.Status(400).JSON(utils.RequestValueRequired("user name"))
	}

	if !utils.IsValidUsername(data["user_name"]) {
		return c.Status(400).JSON(utils.RequestValueValid("user name"))
	}

	// Check if the user name is already in the database
	if userDatabase.CheckUserNameExist(data["user_name"]) != 0 {
		return c.Status(200).JSON(utils.RequestValueInUse("user name"))
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "This user name has not been used yet",
			"inuse":   false,
		})
}
