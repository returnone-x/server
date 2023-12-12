package main

import (
	"os"
	"returnone/config"
	"returnone/middleware"
	"returnone/routes/auth"
	"time"
	"returnone/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// init config
	config.Connect()
	config.GoogleOauth()
	config.GithubOauth()

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

	// cors middleware setup
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://returnone.tech",
		AllowHeaders:  "Origin, Content-Type, Accept",
	}))
	
	// protection Cross-Site Request Forgery (CSRF) attacks
	// *when test csrf must change the ENV*
	if os.Getenv("ENV") == "production" {
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-Csrf-Token",
			CookieName:     "csrf_",
			CookieSameSite: "Strict",
			Expiration:     15 * time.Minute,
			KeyGenerator:   utils.UUID,
		}))
	}

	api_v1 := app.Group("/api/v1")

	// set auth controller
	auth_group := api_v1.Group("/auth")
	auth.Setup(auth_group)

	// set user controller
	user_group := api_v1.Group("/user")
	// set auth middleware(for check does user have right auth)
	user_group.Use(middleware.VerificationAccessToken())
	user.Setup(user_group)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to returnone backend!")
	})

	// app Listen
	app.Listen(":8080")

}
