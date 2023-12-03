package databaseUser

import db "returnone/config"

func GetUserPassword(email string) (string, string) {

	var id string
	var password string

	sqlString := `SELECT id, password FROM users WHERE email = $1;`
	db.DB.QueryRow(sqlString,email).Scan(&id, &password)

	return id, password
}