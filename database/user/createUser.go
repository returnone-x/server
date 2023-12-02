package DatabaseUser

import (
	"returnone/config"
	"returnone/models/user"
	Generate "returnone/utils/generate"
	"time"
)

func CreateUser(
	email string,
	hash_password string,
	user_name string,
) userModles.UserAccount {
	db := config.Connect()

	defer db.Close()

	
	now_time := time.Now()
	sqlString := `
	INSERT INTO users 
	(id, email, password, user_name, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6)
	`

	user_id := Generate.GenerateUserAccountId()

	db.Exec(
		sqlString,
		user_id, email,
		hash_password,
		user_name,
		now_time,
		now_time)
	
	insert_data := userModles.UserAccount{
		Id:             user_id,
		Email:          email,
		Phone:          "",
		Phone_country:  "",
		Password:       hash_password,
		Email_verify:   false,
		Phone_verify:   false,
		Avatar:         "",
		User_name:      user_name,
		Github_connect: "",
		Google_connect: "",
		Create_at:      now_time,
		Update_at:      now_time}

	return insert_data
}
