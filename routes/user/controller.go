package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/middleware"
	"github.com/returnone-x/server/routes/user/question"
	"github.com/returnone-x/server/routes/user/setting"
)

func Setup(app fiber.Router) {
	app.Get("/avatar/:user_id", GetAvatar)

	// add group
	question_group := app.Group("/question")
	question_group.Use(middleware.VerificationAccessToken())
	question.Setup(question_group)

	setting_group := app.Group("/setting")
	setting_group.Use(middleware.VerificationAccessToken())
	setting.Setup(setting_group)
}
