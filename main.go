package main

import (
	"os"
	"returnone/config"
	"returnone/routes/auth"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	db.Connect()
	
	app := fiber.New()

	// Set logger
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} - | ${status} |${latency}     |   ${method}   | ${path} \n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "local",
	}))

	// encrypt cookie
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("ENCRYPT_COOKIE_SECRET"),
	}))
	
	// protection Cross-Site Request Forgery (CSRF) attacks
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     72 * time.Hour,
		KeyGenerator:   utils.UUID,
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
