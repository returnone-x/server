package questionChat

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/middleware"
)

func Setup(app fiber.Router) {
	go runHub()
	app.Get("/questions/:question_id", middleware.VerificationAccessTokenWithoutError(), websocket.New(QuestionsChat))
	app.Post("/questions/new/:question_id", middleware.VerificationAccessToken(), NewMessage)
	app.Delete("/questions/delete/:question_id/:message_id", middleware.VerificationAccessToken(), DeleteMessage)
	app.Post("/questions/update/:question_id/:message_id", middleware.VerificationAccessToken(), UpdateMessage)
	app.Get("/questions/get/:question_id", GetMessage)
}
