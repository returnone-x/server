package auth

import (
	"returnone/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {

	// email signup
	app.Post("/signup", SignUp)
	// email login
	app.Post("/login", LogIn)
	// check does the email or username has been used already
	app.Post("/emailexist", EmailExist)
	app.Post("/usernameexist", UserNameExist)

	// oauth2.0 login or register
	app.Get("/oauth/google", GoogleLogin)
	app.Get("/oauth/callback/google", GoogleCallBack)
	app.Get("/oauth/github", GithubLogin)
	app.Get("/oauth/callback/github", GithubCallBack)

	//check user authorizationa
	app.Get("/authorizationa", middleware.VerificationAccessToken(),CheckAuthorizationa)

	// log out
	app.Get("/logout", middleware.VerificationRefreshToken(),LogOut)

	// for refresh access token
	app.Get("/refresh", middleware.VerificationRefreshToken(),RefreshToken)
}
