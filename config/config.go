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
	redirect_url := ApiUrl() + "/auth/oauth/callback/google"
	GoogleConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  redirect_url,
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
	redirect_url := ApiUrl() + "/auth/oauth/callback/github"
	GithubConfig.GithubLoginConfig = oauth2.Config{
		RedirectURL:  redirect_url,
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user", "read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}

	return GithubConfig.GithubLoginConfig
}

func WebsiteUrl() string {
	if os.Getenv("ENV") == "production" {
		return "https://returnone.tech"
	} else {
		return "https://returnone.nightcat.xyz"
	}
}

func ApiUrl() string {
	if os.Getenv("ENV") == "production" {
		return "https://returnone.tech/api/v1"
	} else {
		return "https://returnone.nightcat.xyz/api/v1"
	}
}