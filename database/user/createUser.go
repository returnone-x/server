package userDatabase

import (
	"returnone/config"
	"returnone/models/user"
	utils "returnone/utils"
	"time"
)

func CreateUser(
	email string,
	hash_password string,
	user_name string,
) (userModles.UserAccount, error){
	now_time := time.Now()
	sqlString := `
	INSERT INTO users 
	(id, email, password, user_name, email_verify, phone_verify, default_2fa, email_2fa, phone_2fa, totp_2fa, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	user_id := utils.GenerateUserAccountId()

	_, err := db.DB.Exec(
		sqlString,
		user_id, email,
		hash_password,
		user_name,
		false,
		false,
		1,
		true,
		false,
		false,
		now_time,
		now_time)
	
	insert_data := userModles.UserAccount{
		Id:             user_id,
		Email:          email,
		Phone:          "",
		Phone_country:  "",
		Email_verify:   false,
		Phone_verify:   false,
		Avatar:         "",
		User_name:      user_name,
		Github_connect: "",
		Google_connect: "",
		Email_2fa:      true,
		Phone_2fa:      false,
		Totp_2fa:       false,
		Create_at:      now_time,
		Update_at:      now_time,
	}

	return insert_data, err
}
