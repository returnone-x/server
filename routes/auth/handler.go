package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/returnone-x/server/config"
	"github.com/returnone-x/server/database/redis"
	tokenDatabase "github.com/returnone-x/server/database/tokens"
	"github.com/returnone-x/server/database/user"
	utils "github.com/returnone-x/server/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func generateAccessTokenExp() time.Time {
	return time.Now().UTC().Add(time.Minute * 60)
}

func generateRefreshTokenExp() time.Time {
	return time.Now().UTC().Add(time.Hour * 24 * 30)
}

func SignUp(c *fiber.Ctx) error {
	var data map[string]string
	// get data from body
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	// check if data is valid return error
	request_data_error := utils.RequestDataRequired(data, []string{"email", "password", "user_name"})
	if request_data_error != nil {
		return c.Status(400).JSON(request_data_error)
	}

	if !utils.IsValidUsername(data["user_name"]) {
		return c.Status(400).JSON(utils.RequestValueValid("username"))
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
		return c.Status(400).JSON(utils.RequestValueInUse("username"))
	}

	hash_password, _ := utils.HashPassword(data["password"])

	save_data, save_data_err := userDatabase.CreateUser(data["email"], hash_password, data["user_name"])

	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", save_data_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error creating user", save_data_err))
	}

	// get user data from database (verify this user is log in or sign up)
	access_token, refresh_token, error_message, err := SetLoginCookies(save_data.Id, c)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    access_token,
		Expires:  generateAccessTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	refresh_token_cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		Expires:  generateRefreshTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	user_id_cookie := fiber.Cookie{
		Name:    "user_id",
		Value:   save_data.Id,
		Expires: generateAccessTokenExp(),
		Domain:  os.Getenv("DOMAIN_NAME"),
	}
	c.Cookie(&user_id_cookie)
	c.Cookie(&refresh_token_cookie)
	c.Cookie(&access_token_cookie)
	return c.Status(200).JSON(
		fiber.Map{
			"status":  "success",
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

	// get user data from database (verify this user is log in or sign up)
	access_token, refresh_token, error_message, err := SetLoginCookies(user_data.Id, c)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
	}

	// just for now this is no need(2023/12/13)
	// if user_data.Default_2fa == 3 {
	// 	if data["otp"] == "" {
	// 		return c.Status(403).JSON(utils.ErrorMessage("OTP is required", nil))
	// 	}

	// 	valid := totp.Validate(data["otp"], user_data.Totp)

	// 	if !valid {
	// 		return c.Status(403).JSON(utils.ErrorMessage("OTP is not valid", nil))
	// 	}
	// }
	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    access_token,
		Expires:  generateAccessTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	refresh_token_cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		Expires:  generateRefreshTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	user_id_cookie := fiber.Cookie{
		Name:    "user_id",
		Value:   user_data.Id,
		Expires: generateAccessTokenExp(),
		Domain:  os.Getenv("DOMAIN_NAME"),
	}
	c.Cookie(&user_id_cookie)
	c.Cookie(&refresh_token_cookie)
	c.Cookie(&access_token_cookie)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
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
	save_state_token_err := redisDB.CreateStringData(state_token_key, state_token, time.Minute*15)

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
	save_token, redis_error := redisDB.GetStrigData(result[0])
	if redis_error != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error get state token", redis_error))
	}

	// if this function is done than run this (i dont want to run this is the middle because it will take some time)
	defer redisDB.DeleteStringData(result[0])

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
	account_result, get_user_data_error := userDatabase.GetGoogleAccount(user_data["id"].(string))

	// if cant found this user than this user is login
	if get_user_data_error == nil {
		access_token, refresh_token, error_message, err := SetLoginCookies(account_result.Id, c)

		if err != nil {
			return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
		}

		//set cookies
		access_token_cookie := fiber.Cookie{
			Name:     "accessToken",
			Value:    access_token,
			Expires:  generateAccessTokenExp(),
			HTTPOnly: true,
			Secure:   true,
			Domain:   os.Getenv("DOMAIN_NAME"),
		}
		refresh_token_cookie := fiber.Cookie{
			Name:     "refreshToken",
			Value:    refresh_token,
			Expires:  generateRefreshTokenExp(),
			HTTPOnly: true,
			Secure:   true,
			Domain:   os.Getenv("DOMAIN_NAME"),
		}
		user_id_cookie := fiber.Cookie{
			Name:    "user_id",
			Value:   account_result.Id,
			Expires: generateAccessTokenExp(),
			Domain:  os.Getenv("DOMAIN_NAME"),
		}
		c.Cookie(&user_id_cookie)
		c.Cookie(&refresh_token_cookie)
		c.Cookie(&access_token_cookie)
		return c.Status(200).Redirect(config.WebsiteUrl() + "/logincomplete")
	}

	// check does user already create a account
	_, get_user_data_by_email_error := userDatabase.GetGoogleAccountWithEmail(user_data["email"].(string))
	if get_user_data_by_email_error != nil {
		return c.Status(400).JSON(utils.ErrorMessage("You have already use this google account's email register", nil))
	}

	// if this user didn't sign up than create user data
	save_account_result, save_data_err := userDatabase.CreateUserWithGoogleLogin(user_data["email"].(string), user_data["picture"].(string), user_data["id"].(string))
	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", save_data_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error creating user", save_data_err))
	}

	access_token, refresh_token, error_message, err := SetLoginCookies(save_account_result.Id, c)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    access_token,
		Expires:  generateAccessTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	refresh_token_cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		Expires:  generateRefreshTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	user_id_cookie := fiber.Cookie{
		Name:    "user_id",
		Value:   save_account_result.Id,
		Expires: generateAccessTokenExp(),
		Domain:  os.Getenv("DOMAIN_NAME"),
	}
	c.Cookie(&user_id_cookie)
	c.Cookie(&refresh_token_cookie)
	c.Cookie(&access_token_cookie)
	// for front end to check is this is the sign up if it is than popup a modal to let user change username
	c.Cookie(&fiber.Cookie{
		Name:    "first_login",
		Value:   "1",
		Expires: time.Now().UTC().Add(time.Second * 10),
	})

	return c.Status(200).Redirect(config.WebsiteUrl() + "/logincomplete")
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
	save_state_token_err := redisDB.CreateStringData(state_token_key, state_token, time.Minute*15)

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
	save_token, redis_error := redisDB.GetStrigData(result[0])
	if redis_error != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error get state token", redis_error))
	}

	// if this function is done than run this (i dont want to run this is the middle because it will take some time)
	defer redisDB.DeleteStringData(result[0])

	//check the states
	if result[1] != save_token {
		return c.Status(500).JSON(utils.ErrorMessage("States don't match", nil))
	}

	// so i dont know how does this work
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
	fmt.Println(user_data)

	// get user email
	email_client := &http.Client{}

	email_req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("User data fetch failed", err))
	}

	email_req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	email_req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	email_req.Header.Set("Accept", "application/vnd.github+json")

	// get the user data
	email_resp, err := email_client.Do(email_req)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("User data fetch failed", err))
	}

	user_emails_data_byte, err := io.ReadAll(email_resp.Body)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("JSON parsing failed", err))
	}
	emails_json_str := string(user_emails_data_byte)
	var user_emails []map[string]interface{}
	err = json.Unmarshal([]byte(emails_json_str), &user_emails)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("JSON unmarshal failed", err))
	}

	var primary_email string

	// get primary email
	for _, userMap := range user_emails {
		if primary, ok := userMap["primary"].(bool); ok && primary {

			primary_email, ok = userMap["email"].(string)
			if !ok {
				return c.Status(500).JSON(utils.ErrorMessage("JSON unmarshal failed", err))
			}
		}
	}

	// convert user data id to string
	user_github_id := fmt.Sprintf("%f", user_data["id"])

	// get user data from database (verify this user is log in or sign up)
	account_reslut, get_user_data_error := userDatabase.GetGithubAccount(user_github_id)

	// if cant found this user than this user is login
	if get_user_data_error == nil {
		// get user data from database (verify this user is log in or sign up)
		access_token, refresh_token, error_message, err := SetLoginCookies(account_reslut.Id, c)

		if err != nil {
			return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
		}

		//set cookies
		access_token_cookie := fiber.Cookie{
			Name:     "accessToken",
			Value:    access_token,
			Expires:  generateAccessTokenExp(),
			HTTPOnly: true,
			Secure:   true,
			Domain:   os.Getenv("DOMAIN_NAME"),
		}
		refresh_token_cookie := fiber.Cookie{
			Name:     "refreshToken",
			Value:    refresh_token,
			Expires:  generateRefreshTokenExp(),
			HTTPOnly: true,
			Secure:   true,
			Domain:   os.Getenv("DOMAIN_NAME"),
		}
		user_id_cookie := fiber.Cookie{
			Name:    "user_id",
			Value:   account_reslut.Id,
			Expires: generateAccessTokenExp(),
			Domain:  os.Getenv("DOMAIN_NAME"),
		}
		c.Cookie(&user_id_cookie)
		c.Cookie(&refresh_token_cookie)
		c.Cookie(&access_token_cookie)
		return c.Status(200).Redirect(config.WebsiteUrl() + "/logincomplete")
	}

	// if this user didn't sign up than create user data
	save_account_reslut, save_data_err := userDatabase.CreateUserWithGithubLogin(primary_email, user_github_id, user_data["avatar_url"].(string))

	// get user data from database (verify this user is log in or sign up)
	access_token, refresh_token, error_message, err := SetLoginCookies(save_account_reslut.Id, c)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    access_token,
		Expires:  generateAccessTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	refresh_token_cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		Expires:  generateRefreshTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	user_id_cookie := fiber.Cookie{
		Name:    "user_id",
		Value:   save_account_reslut.Id,
		Expires: generateAccessTokenExp(),
		Domain:  os.Getenv("DOMAIN_NAME"),
	}
	c.Cookie(&user_id_cookie)
	c.Cookie(&refresh_token_cookie)
	c.Cookie(&access_token_cookie)
	// for front end to check is this is the sign up if it is than popup a modal to let user change username
	c.Cookie(&fiber.Cookie{
		Name:    "first_login",
		Value:   "1",
		Expires: time.Now().UTC().Add(time.Second * 10),
	})
	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", user_data, "| Message:", save_data_err)
		return c.Status(500).JSON(utils.ErrorMessage("Error creating user", save_data_err))
	}

	c.Cookie(&refresh_token_cookie)
	c.Cookie(&access_token_cookie)

	return c.Status(200).Redirect(config.WebsiteUrl() + "/logincomplete")
}

func RefreshToken(c *fiber.Ctx) error {
	// Get the refresh token details from the local
	// This will be set when through the middleware
	token := c.Locals("refresh_token_context").(*jwt.Token)

	// Change token to jwt mapclaims for later access data
	claims := token.Claims.(jwt.MapClaims)

	// get token id from claims
	token_id := claims["token_id"].(string)

	// this is just for convert type to the int
	ot := claims["used_time"].(string)
	old_used_times, _ := strconv.Atoi(ot)

	// get data from database
	refresh_token_context, err := tokenDatabase.GetTokenData(token_id)

	// if got any error
	if err != nil {
		return c.Status(401).JSON(utils.ErrorMessage("Error get refresh token detils", err))
	}

	// check the used times is the same (if not the same this account is been hacked)
	if old_used_times != refresh_token_context.Used_time {
		// so delete the token for sure the hacker cant use this session anymore
		tokenDatabase.DeleteToken(token_id)
		return c.Status(403).JSON(
			fiber.Map{
				"status":  "error",
				"message": "This refresh token has been used",
			})
	}

	// add the number of uses since login
	new_used_times := old_used_times
	new_used_times++

	// get user id
	user_id := claims["user_id"].(string)

	// create new access token and refresh token (also update the refresh token used time)
	access_token, refresh_token, error_message, err := SetRefreshCookies(user_id, token_id, new_used_times, c)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage(error_message, err))
	}

	//set cookies
	access_token_cookie := fiber.Cookie{
		Name:     "accessToken",
		Value:    access_token,
		Expires:  generateAccessTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	refresh_token_cookie := fiber.Cookie{
		Name:     "refreshToken",
		Value:    refresh_token,
		Expires:  generateRefreshTokenExp(),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	}
	user_id_cookie := fiber.Cookie{
		Name:    "user_id",
		Value:   user_id,
		Expires: generateAccessTokenExp(),
		Domain:  os.Getenv("DOMAIN_NAME"),
	}
	c.Cookie(&user_id_cookie)
	c.Cookie(&refresh_token_cookie)
	c.Cookie(&access_token_cookie)

	return c.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Successfully refresh token",
		})
}

func LogOut(c *fiber.Ctx) error {
	token := c.Locals("refresh_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	token_id := claims["token_id"].(string)
	result, err := tokenDatabase.DeleteToken(token_id)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When delete this session database got some error", err))
	}

	result.RowsAffected()

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    "None",
		Expires:  time.Now().UTC().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "None",
		Expires:  time.Now().UTC().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		Secure:   true,
		Domain:   os.Getenv("DOMAIN_NAME"),
	})

	c.Cookie(&fiber.Cookie{
		Name:    "user_id",
		Value:   "None",
		Expires: time.Now().UTC().Add(-(time.Hour * 2)),
		Domain:  os.Getenv("DOMAIN_NAME"),
	})
	return c.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Successfully log out",
		})
}

func CheckAuthorizationa(c *fiber.Ctx) error {
	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	if claims["exp"].(float64)-float64(time.Now().UTC().Unix()) < 60*20 {
		return c.Status(200).JSON(
			fiber.Map{
				"status":  "success",
				"message": "The token will expire in the near future",
			})
	}
	// cuz middleware has already checked the authorization so dont need to check again
	return c.Status(200).JSON(
		fiber.Map{
			"status":  "success",
			"message": "This user has been authorized",
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
			"status":  "success",
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
			"status":  "success",
			"message": "This user name has not been used yet",
			"inuse":   false,
		})
}
