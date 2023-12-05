package auth

import (
	"returnone/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {

	// email signup
	app.Post("/signup", SignUp)
	// email login
	app.Get("/login", LogIn)
	// check does the email or username has been used already
	app.Get("/emailexist", EmailExist)
	app.Get("/usernameexist", UserNameExist)

	// oauth2.0 login or register
	app.Get("/oauth/google", GoogleLogin)
	app.Get("/oauth/callback/google", GoogleCallBack)
	app.Get("/oauth/github", GithubLogin)
	app.Get("/oauth/callback/github", GithubCallBack)

	//check user authorizationa
	app.Get("/authorizationa", middleware.VerificationAccessToken(),CheckAuthorizationa)

	// for refresh access token
	app.Get("/refresh", middleware.VerificationRefreshToken(),RefreshToken)
}
