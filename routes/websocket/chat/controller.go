package questionChat

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/middleware"
)

func Setup(app fiber.Router) {
	app.Get("/questions/:questionId", middleware.VerificationAccessTokenWithoutError(), websocket.New(QuestionsChat))
}
