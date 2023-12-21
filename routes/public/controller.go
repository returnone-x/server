package public

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Get("/question/:id", GetQuestion)
}
