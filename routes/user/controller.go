package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/middleware"
	"github.com/returnone-x/server/routes/user/question"
)

func Setup(app fiber.Router) {
	app.Post("/rename", middleware.VerificationAccessToken(), Rename)
	app.Get("/avatar", middleware.VerificationRefreshToken(), GetAvatar)

	// add group
	auth_group := app.Group("/question")
	auth_group.Use(middleware.VerificationAccessToken())
	question.Setup(auth_group)
}
