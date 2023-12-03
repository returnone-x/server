package userModles

import "time"

type UserAccount struct {
	Id             string    `json:"id"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Phone_country  string    `json:"phone_country"`
	Password       string    `json:"password"`
	Avatar         string    `json:"avatar"`
	User_name      string    `json:"user_name"`
	Email_verify   bool      `json:"email_verify"`
	Phone_verify   bool      `json:"phone_verify"`
	Github_connect string    `json:"github_connect"`
	Google_connect string    `json:"google_connect"`
	Phone_2fa      bool      `json:"phone_2fa"`
	Email_2fa      bool      `json:"email_2fa"`
	Totp_2fa       bool      `json:"totp_2fa"`
	Totp           string    `json:"totp"`
	Default_2fa    int       `json:"default_2fa"`
	Create_at      time.Time `json:"create_at"`
	Update_at      time.Time `json:"update_at"`
}

// table name: users
