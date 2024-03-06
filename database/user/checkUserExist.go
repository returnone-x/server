package userDatabase

import (
	"github.com/returnone-x/server/config"
)

func CheckUserEmailExist(email string) int {

	var count int

	sqlString := `SELECT COUNT(*) FROM users WHERE email = $1;`
	config.DB.QueryRow(sqlString, email).Scan(&count)

	return count
}

func CheckUserNameExist(username string) int {

	var count int

	sqlString := `SELECT COUNT(*) FROM users WHERE username = $1;`
	config.DB.QueryRow(sqlString, username).Scan(&count)

	return count
}

func CheckUserGoogleAccountExist(username string) int {

	var count int

	sqlString := `SELECT COUNT(*) FROM users WHERE google_connect = $1;`
	config.DB.QueryRow(sqlString, username).Scan(&count)

	return count
}
