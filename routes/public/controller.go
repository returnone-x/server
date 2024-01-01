package public

import (
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/middleware"
)

func Setup(app fiber.Router) {
	app.Get("/question/:id", middleware.VerificationAccessTokenWithoutError(), GetQuestion)
	app.Get("/question/comments/:question_id", GetQuestionComment)
}
