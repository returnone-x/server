package userDatabase

import (
	db "github.com/returnone-x/server/config"
)

type UserAuthData struct {
	Id           string `json:"id"`
	Password     string `json:"password"`
	Email_verify bool   `json:"email_verify"`
	Email_2fa    bool   `json:"email_2fa"`
	Phone_2fa    bool   `json:"phone_2fa"`
	Totp_2fa     bool   `json:"totp_2fa"`
	Totp         string `json:"totp"`
	Default_2fa  int    `json:"default_2fa"`
}

type OauthData struct {
	Id           string `json:"id"`
	Email_verify bool   `json:"email_verify"`
	Email_2fa    bool   `json:"email_2fa"`
	Phone_2fa    bool   `json:"phone_2fa"`
	Totp_2fa     bool   `json:"totp_2fa"`
	Totp         string `json:"totp"`
	Default_2fa  int    `json:"default_2fa"`
}

func GetUserPassword(email string) (UserAuthData, error) {

	var user_data UserAuthData

	sqlString := `SELECT id, password, email_verify, email_2fa, phone_2fa, totp_2fa, totp, default_2fa FROM users WHERE email = $1;`
	err := db.DB.QueryRow(sqlString, email).Scan(
		&user_data.Id,
		&user_data.Password,
		&user_data.Email_verify,
		&user_data.Email_2fa,
		&user_data.Phone_2fa,
		&user_data.Totp_2fa,
		&user_data.Totp,
		&user_data.Default_2fa)
	return user_data, err
}

func GetGoogleAccount(google_id string) (OauthData, error) {

	var user_data OauthData

	sqlString := `SELECT id, email_verify, email_2fa, phone_2fa, totp_2fa, totp, default_2fa FROM users WHERE google_connect = $1;`
	err := db.DB.QueryRow(sqlString, google_id).Scan(
		&user_data.Id,
		&user_data.Email_verify,
		&user_data.Email_2fa,
		&user_data.Phone_2fa,
		&user_data.Totp_2fa,
		&user_data.Totp,
		&user_data.Default_2fa,
			
	)
	return user_data, err
}

func GetGithubAccount(github_id string) (OauthData, error) {

	var user_data OauthData

	sqlString := `SELECT id, email_verify, email_2fa, phone_2fa, totp_2fa, totp, default_2fa FROM users WHERE github_connect = $1;`
	err := db.DB.QueryRow(sqlString, github_id).Scan(
		&user_data.Id,
		&user_data.Email_verify,
		&user_data.Email_2fa,
		&user_data.Phone_2fa,
		&user_data.Totp_2fa,
		&user_data.Totp,
		&user_data.Default_2fa)
	return user_data, err
}

