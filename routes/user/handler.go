package user

import "github.com/gofiber/fiber/v2"

func Rename(c *fiber.Ctx) error {
	
	return c.SendStatus(200)
}