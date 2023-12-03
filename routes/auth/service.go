package auth

import (
	"fmt"
	"log"
	"returnone/database/user"
	"returnone/utils/crypto"
	"returnone/utils/sendError"
	"returnone/utils/valid"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
	request_data_error := sendError.RequestDataError(data, []string{"email", "password", "user_name"})
	if request_data_error != "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": fmt.Sprintf("%v is required", request_data_error),
			})
	}

	if !valid.IsValidUsername(data["user_name"]) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "This user name is not valid",
		})
	}

	if !valid.IsValidEmail(data["email"]) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "This email address is not valid",
		})
	}

	// check the email has already been used
	if databaseUser.CheckUserEmailExist(data["email"]) != 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "This email is already in use",
			})
	}

	// check the user name has already been used
	if databaseUser.CheckUserNameExist(data["user_name"]) != 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "This user name is already in use",
			})
	}

	hash_password, _ := crypto.HashPassword(data["password"])

	save_data, save_data_err := databaseUser.CreateUser(data["email"], hash_password, data["user_name"])

	if save_data_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", save_data_err)
		c.Status(fiber.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Could not create this user when inserting into database",
			"error":   save_data_err,
		})
	}

	//24 hours later time(for access token)
	twenty_four_hours_later := time.Now().Add(time.Hour * 24)
	//60 days later time(for refresh tokne)
	sixty_days_later := time.Now().Add(time.Hour * 24 * 60)

	//generate Jwt token
	access_token, access_token_err := crypto.GenerateJwtToken(save_data.Id, "accessToken", twenty_four_hours_later.Unix())
	refresh_token, refresh_token_err := crypto.GenerateJwtToken(save_data.Id, "refreshToken", sixty_days_later.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", refresh_token_err)
		c.Status(fiber.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cloud not login when generating refresh token",
			"error":   access_token_err,
		})
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", access_token_err)
		c.Status(fiber.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cloud not login when generating access token",
			"error":   access_token_err,
		})
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
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid post request",
			})
	}

	if !valid.IsValidEmail(data["email"]) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "This email address is not valid",
		})
	}

	// check if data is valid return error
	request_data_error := sendError.RequestDataError(data, []string{"email", "password"})
	if request_data_error != "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": fmt.Sprintf("%v is required", request_data_error),
			})
	}

	// check the email has already been used
	if databaseUser.CheckUserEmailExist(data["email"]) == 0 {
		return c.Status(401).JSON(
			fiber.Map{
				"success": false,
				"message": "Password or email is incorrect",
			})
	}

	//Get hash password from database
	user_id, hash_password := databaseUser.GetUserPassword(data["email"])

	//check password is correct
	check_password := crypto.CheckPasswordHash(data["password"], hash_password)

	// if password not correct
	if !check_password {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Password or email is incorrect",
		})
	}

	//24 hours later time(for access token)
	twenty_four_hours_later := time.Now().Add(time.Hour * 24)
	//60 days later time(for refresh tokne)
	sixty_days_later := time.Now().Add(time.Hour * 24 * 60)

	//generate Jwt token
	access_token, access_token_err := crypto.GenerateJwtToken(user_id, "accessToken", twenty_four_hours_later.Unix())
	refresh_token, refresh_token_err := crypto.GenerateJwtToken(user_id, "refreshToken", sixty_days_later.Unix())

	//handle errors
	if refresh_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", refresh_token_err)
		c.Status(fiber.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cloud not login when generating refresh token",
			"error":   access_token_err,
		})
	}
	if access_token_err != nil {
		log.Println("| Path:", c.Path(), "| Data:", data, "| Message:", access_token_err)
		c.Status(fiber.StatusInternalServerError)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Cloud not login when generating access token",
			"error":   access_token_err,
		})
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

func EmailExist(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid post request",
			})
	}

	if data["email"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Email is required",
			})
	}

	if !valid.IsValidEmail(data["email"]) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "This email address is not valid",
		})
	}

	// Check if the email is already in the database
	if databaseUser.CheckUserEmailExist(data["email"]) != 0 {
		return c.Status(200).JSON(
			fiber.Map{
				"success": true,
				"message": "This email is already in use",
				"inuse":   true,
			})
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
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Invalid post request",
			})
	}

	if data["user_name"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "User name is required",
			})
	}

	if !valid.IsValidUsername(data["user_name"]) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "This user name is not valid",
		})
	}

	// Check if the user name is already in the database
	if databaseUser.CheckUserNameExist(data["user_name"]) != 0 {
		return c.Status(200).JSON(
			fiber.Map{
				"success": true,
				"message": "This user name is already in use",
				"inuse":   true,
			})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "This user name has not been used yet",
			"inuse":   false,
		})
}
