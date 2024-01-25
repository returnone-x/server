package userDatabase

import (
	db "github.com/returnone-x/server/config"
	"github.com/returnone-x/server/models/user"
)

func GetUserDetil(id string) (userModles.UserAccount, error) {
	var user_data userModles.UserAccount

	sqlString := `SELECT 
		id,
		email,
		COALESCE(phone, '') AS phone,
		COALESCE(phone_country, '') AS phone_country,
		COALESCE(password, '') AS password,
		email_verify,
		phone_verify,
		avatar,
		COALESCE(display_name, '') AS display_name,
		user_name,
		COALESCE(github_connect, '') AS github_connect,
		COALESCE(google_connect, '') AS google_connect,
		email_2fa,
		phone_2fa,
		totp_2fa,
		COALESCE(totp, '') AS totp,
		default_2fa,
		create_at,
		update_at
	FROM users WHERE id = $1;`
	err := db.DB.QueryRow(sqlString, id).Scan(
		&user_data.Id,
		&user_data.Email,
		&user_data.Phone,
		&user_data.Phone_country,
		&user_data.Password,
		&user_data.Email_verify,
		&user_data.Phone_verify,
		&user_data.Avatar,
		&user_data.Display_name,
		&user_data.User_name,
		&user_data.Github_connect,
		&user_data.Google_connect,
		&user_data.Phone_2fa,
		&user_data.Email_2fa,
		&user_data.Totp_2fa,
		&user_data.Totp,
		&user_data.Default_2fa,
		&user_data.Create_at,
		&user_data.Update_at,
	)

	return user_data, err
}
