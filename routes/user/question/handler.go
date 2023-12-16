package question

import (
	"github.com/gofiber/fiber/v2"
	utils "github.com/returnone-x/server/utils"
)

func NewPost(c *fiber.Ctx) error {
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

	request_data_error := utils.RequestDataRequired(data, []string{"email", "password", "user_name"})
	if request_data_error != nil {
		return c.Status(400).JSON(request_data_error)
	}
	return c.SendStatus(200)
}
