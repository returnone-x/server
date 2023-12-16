package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/routes/user/question"
)

func Setup(app fiber.Router) {
	app.Post("/rename", Rename)
	app.Get("/avatar", GetAvatar)
	
	// add group
	auth_group := app.Group("/question")
	question.Setup(auth_group)
}
