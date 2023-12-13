package userDatabase

import db "returnone/config"

func GetUserAvatar(id string) (string, error) {

	var avatar string

	sqlString := `SELECT avatar FROM users WHERE id = $1;`
	err := db.DB.QueryRow(sqlString, id).Scan(
		&avatar)

	return avatar, err
}
