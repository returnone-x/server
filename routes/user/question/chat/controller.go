package questionChat

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {

	app.Post("/message/new/:question_id", NewMessage)
}