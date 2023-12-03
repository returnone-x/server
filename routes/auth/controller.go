package auth

import "github.com/gofiber/fiber/v2"

func Setup(app fiber.Router) {

	app.Post("/signup", SignUp)
	app.Get("/login", LogIn)
	// app.Post("/logout", LogOut)
	app.Get("/emailexist", EmailExist)
	app.Get("/usernameexist", UserNameExist)

	// oauth2.0 login or register
	app.Get("/oauth/google", GoogleLogin)
	app.Get("/oauth/callback/google", GoogleCallBack)
	// app.Post("/oauth/github")
	// app.Post("/oauth/callback/github")
}
