package main

import (
	"returnone/routes/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Set logger
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} - | ${status} |${latency}     |   ${method}   | ${path} \n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "local",
	}))

	api_v1 := app.Group("/v1")

	// set auth controller
	auth_group := api_v1.Group("/auth")
	auth.Setup(auth_group)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to returnone backend!")
	})

	// app Listen
	app.Listen(":8080")

}
