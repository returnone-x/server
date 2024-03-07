package userDatabase

import (
	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	userModles "github.com/returnone-x/server/models/user"
)

func GetUserDetil(id string) (userModles.UserDetil, error) {
	var user_data userModles.UserDetil

	sqlString := `
	SELECT 
    u.id,
    u.email,
    COALESCE(u.phone, '') AS phone,
    COALESCE(u.phone_country, '') AS phone_country,
    COALESCE(u.password, '') AS password,
    u.email_verify,
    u.phone_verify,
    u.avatar,
    COALESCE(u.display_name, '') AS display_name,
    u.username,
    COALESCE(u.github_connect, '') AS github_connect,
    COALESCE(u.google_connect, '') AS google_connect,
    u.email_2fa,
    u.phone_2fa,
    u.totp_2fa,
    COALESCE(u.totp, '') AS totp,
    u.default_2fa,
    up.bio,
    up.public_email,
    up.pronouns,
    up.related_links,
    u.create_at,
    u.update_at
FROM
    users u
LEFT JOIN
    user_profile up ON up.id = u.id 
WHERE 
    u.id = $1;
`
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
		&user_data.Username,
		&user_data.Github_connect,
		&user_data.Google_connect,
		&user_data.Phone_2fa,
		&user_data.Email_2fa,
		&user_data.Totp_2fa,
		&user_data.Totp,
		&user_data.Default_2fa,
		&user_data.Bio,
		&user_data.Public_email,
		&user_data.Pronouns,
		(*pq.StringArray)(&user_data.Related_links),
		&user_data.Create_at,
		&user_data.Update_at,
	)

	return user_data, err
}
