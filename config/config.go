package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() oauth2.Config {

	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://127.0.0.1:8080/v1/auth/oauth/callback/google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}
