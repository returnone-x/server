package user

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	userDatabase "github.com/returnone-x/server/database/user"
	utils "github.com/returnone-x/server/utils"
)

func GetAvatar(c *fiber.Ctx) error {
	params := c.AllParams()

	if params["user_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}
	user_id := params["user_id"]
	avatar, err := userDatabase.GetUserAvatar(user_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(utils.ErrorMessage("Can't find this user's avatar", err))
		} else {
			return c.Status(500).JSON(utils.ErrorMessage("Error when get user's avatar from database", err))
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful get user avatar",
		"data":    avatar})
}
