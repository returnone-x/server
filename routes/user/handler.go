package user

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	userDatabase "github.com/returnone-x/server/database/user"
	utils "github.com/returnone-x/server/utils"
)

func Rename(c *fiber.Ctx) error {

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

	// check is this username valid
	if data["new_username"] == "" {
		return c.Status(400).JSON(utils.RequestValueRequired("new username"))
	}

	if !utils.IsValidUsername(data["new_username"]) {
		return c.Status(400).JSON(utils.RequestValueValid("new username"))
	}

	// check the user name has already been used
	if userDatabase.CheckUserNameExist(data["new_username"]) != 0 {
		return c.Status(400).JSON(utils.RequestValueInUse("new username"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	result, err := userDatabase.Rename(user_id, data["new_username"])

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When update database got some error", err))
	}

	affected_row, _ := result.RowsAffected()
	// if didnt update anything than throw a error
	if affected_row == 0 {
		return c.Status(400).JSON(utils.ErrorMessage("We cant find this user in database", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully update username",
	})
}

func GetAvatar(c *fiber.Ctx) error {
	token := c.Locals("refresh_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	user_id := claims["user_id"].(string)
	avatar, err := userDatabase.GetUserAvatar(user_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(utils.ErrorMessage("Can't find this user's avatar", err))
		} else {
			return c.Status(500).JSON(utils.ErrorMessage("Error when get user's avatar from database", err))
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful get user avatar",
		"data":    avatar})
}
