package question

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Post("/new" , NewPost)
}