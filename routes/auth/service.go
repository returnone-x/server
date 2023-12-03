package auth

import (
	"fmt"
	"log"
	"returnone/database/user"
	utils "returnone/utils"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
)


//24 hours later time(for access token)
var twenty_four_hours_later = time.Now().Add(time.Hour * 24)
//60 days later time(for refresh tokne)
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

	// check the email has already been used
	if userDatabase.CheckUserEmailExist(data["email"]) == 0 {
		return c.Status(401).JSON(utils.RequestValueValid("password or email"))
	}

	//Get hash password from database
	user_data := userDatabase.GetUserPassword(data["email"])

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
		if data["otp"] == ""{
			return c.Status(403).JSON(utils.ErrorMessage("OTP is required", nil))
		}
		fmt.Println(user_data.Totp)

		valid := totp.Validate(data["otp"], user_data.Totp)
	
		fmt.Println("current one-time password is:", user_data.Totp)
	
		fmt.Println("verify OTP success:", valid)
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
	return c.SendString("test")
}

func GoogleCallBack(c *fiber.Ctx) error {
	return c.SendString("test")
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
