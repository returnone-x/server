package auth

import "github.com/gofiber/fiber/v2"

func Setup(app fiber.Router) {
	
	app.Post("/signup", SignUp)
	// app.Post("/login", LogIn)
	// app.Post("/logout", LogOut)
	app.Get("/emailexist", EmailExist)
	app.Get("/usernameexist", UserNameExist)
}