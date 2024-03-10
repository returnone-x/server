package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/routes/auth"
	"github.com/returnone-x/server/routes/public"
	"github.com/returnone-x/server/routes/user"
	"github.com/returnone-x/server/routes/websocket"
)

func routes(app *fiber.App){
	api_v1 := app.Group("/api/v1")

	// set auth controller
	auth_group := api_v1.Group("/auth")
	auth.Setup(auth_group)

	// set public controller(most for get public resource)
	public_group := api_v1.Group("/public")
	public.Setup(public_group)

	// set user controller
	user_group := api_v1.Group("/user")
	user.Setup(user_group)

	// set websocket controller
	websocket_group := api_v1.Group("/ws")
	chat.Setup(websocket_group)
}