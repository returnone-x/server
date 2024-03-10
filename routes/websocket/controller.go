package chat

import (
	"github.com/gofiber/fiber/v2"
	chat "github.com/returnone-x/server/routes/websocket/chat"
)

func Setup(app fiber.Router) {
	
	chat_group := app.Group("/chat")
	chat.Setup(chat_group)
}