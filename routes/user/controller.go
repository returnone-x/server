package user

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Post("/rename" , Rename)
}
