package DatabaseUser

import (
	"returnone/config"
)

func CheckUserEmailExist(email string) int {
	db := config.Connect()

	defer db.Close()

	var count int

	sqlString := `SELECT COUNT(*) FROM users WHERE email = $1;`
	db.QueryRow(sqlString,email).Scan(&count)

	return count
}

func CheckUserNameExist(user_name string) int {
	db := config.Connect()

	defer db.Close()

	var count int

	sqlString := `SELECT COUNT(*) FROM users WHERE user_name = $1;`
	db.QueryRow(sqlString,user_name).Scan(&count)

	return count
}