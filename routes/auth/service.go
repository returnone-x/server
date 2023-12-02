package auth

import (
	"fmt"
	"returnone/database/user"
	"returnone/utils/crypto"
	SendError "returnone/utils/sendError"

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
	request_data_error := SendError.RequestDataError(data, []string{"email", "password", "user_name"})
	if request_data_error != ""{
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": fmt.Sprintf("%v is required", request_data_error),
			})
	}

	// check the email has already been used
	if DatabaseUser.CheckUserEmailExist(data["email"]) != 0{
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "This email is already in use",				
			})
	}

	// check the user name has already been used
	if DatabaseUser.CheckUserNameExist(data["user_name"]) != 0{
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "This user name is already in use",				
			})
	}

	hash_password, _ := crypto.HashPassword(data["password"])
	
	save_data := DatabaseUser.CreateUser(data["email"], hash_password, data["user_name"])

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "Sign up successfully",
			"data": save_data,
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

	if data["email"] == ""{
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "Email is required",
			})
	}

	// Check if the email is already in the database
	if DatabaseUser.CheckUserEmailExist(data["email"]) != 0{
		return c.Status(200).JSON(
			fiber.Map{
				"success": true,
				"message": "This email is already in use",		
				"inuse": true,		
			})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "This email has not been used yet",
			"inuse": false,
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

	if data["user_name"] == ""{
		return c.Status(400).JSON(
			fiber.Map{
				"success": false,
				"message": "User name is required",
			})
	}

	// Check if the user name is already in the database
	if DatabaseUser.CheckUserNameExist(data["user_name"]) != 0{
		return c.Status(200).JSON(
			fiber.Map{
				"success": true,
				"message": "This user name is already in use",
				"inuse": true,				
			})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success": true,
			"message": "This user name has not been used yet",
			"inuse": false,
		})
}
