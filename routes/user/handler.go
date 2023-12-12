package user

import (
	userDatabase "returnone/database/user"
	utils "returnone/utils"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func Rename(c *fiber.Ctx) error {

	var data map[string]string
	// get data from body
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "success",
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
	user_id := claims["user_id"].(string)

	result, err := userDatabase.Rename(user_id, data["new_username"])

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When update database got some error", err))
	}

	affected_row, _ := result.RowsAffected()

	if affected_row == 0 {
		return c.Status(400).JSON(utils.ErrorMessage("We cant find this user in database", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully update username",
	})
}
