package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type GoogleConfigType struct {
	GoogleLoginConfig oauth2.Config
}

var GoogleConfig GoogleConfigType

func GoogleOauth() oauth2.Config {

	GoogleConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  "http://127.0.0.1:8080/v1/auth/oauth/callback/google",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return GoogleConfig.GoogleLoginConfig
}

type GithubConfigType struct {
	GithubLoginConfig oauth2.Config
}

var GithubConfig GithubConfigType

func GithubOauth() oauth2.Config {

	GithubConfig.GithubLoginConfig = oauth2.Config{
		RedirectURL:  "http://127.0.0.1:8080/v1/auth/oauth/callback/github",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes: []string{"user", "read:user"},
		Endpoint: github.Endpoint,
	}

	return GithubConfig.GithubLoginConfig
}
