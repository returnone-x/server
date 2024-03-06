package userDatabase

import (
	"fmt"
	"time"

	"github.com/returnone-x/server/config"
	userModles "github.com/returnone-x/server/models/user"
	utils "github.com/returnone-x/server/utils"
)

func CreateUser(
	email string,
	hash_password string,
	username string,
) (userModles.UserAccount, error) {
	now_time := time.Now().UTC()
	sqlString := `
	INSERT INTO users 
	(id, email, password, username, email_verify, phone_verify, default_2fa, email_2fa, phone_2fa, totp_2fa, totp, create_at, update_at, avatar) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	user_id := utils.GenerateUserAccountId()
	// the avatar just for now(feture will change to random avatar)
	_, err := config.DB.Exec(
		sqlString,
		user_id, email,
		hash_password,
		username,
		false,
		false,
		1,
		true,
		false,
		false,
		"",
		now_time,
		now_time,
		"https://i1.sndcdn.com/artworks-DMKEsjVymB5A2teD-yr6bng-t240x240.jpg")

	insert_data := userModles.UserAccount{
		Id:             user_id,
		Email:          email,
		Phone:          "",
		Phone_country:  "",
		Email_verify:   false,
		Phone_verify:   false,
		Avatar:         "",
		Username:       username,
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

func CreateUserWithGoogleLogin(
	email string,
	avatar string,
	google_id string,
) (userModles.UserAccount, error) {
	now_time := time.Now().UTC()
	sqlString := `
	INSERT INTO users 
	(id, email, username, email_verify, phone_verify, default_2fa, email_2fa, phone_2fa, totp_2fa, totp, avatar, google_connect, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	user_id := utils.GenerateUserAccountId()

	_, err := config.DB.Exec(
		sqlString,
		user_id, email,
		user_id,
		false,
		false,
		0,
		true,
		false,
		false,
		"",
		avatar,
		google_id,
		now_time,
		now_time)

	insert_data := userModles.UserAccount{
		Id:             user_id,
		Email:          email,
		Phone:          "",
		Phone_country:  "",
		Email_verify:   true,
		Phone_verify:   false,
		Avatar:         avatar,
		Username:       user_id,
		Github_connect: "",
		Google_connect: google_id,
		Email_2fa:      true,
		Phone_2fa:      false,
		Totp_2fa:       false,
		Create_at:      now_time,
		Update_at:      now_time,
	}
	fmt.Sprintln(err)
	return insert_data, err
}

func CreateUserWithGithubLogin(
	github_email string,
	github_id string,
	avatar string,
) (userModles.UserAccount, error) {
	now_time := time.Now().UTC()
	sqlString := `
	INSERT INTO users 
	(id, email, username, email_verify, phone_verify, default_2fa, email_2fa, phone_2fa, totp_2fa, totp, avatar, github_connect, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	user_id := utils.GenerateUserAccountId()

	_, err := config.DB.Exec(
		sqlString,
		user_id, github_email,
		user_id,
		false,
		false,
		0,
		true,
		false,
		false,
		"",
		avatar,
		github_id,
		now_time,
		now_time)

	insert_data := userModles.UserAccount{
		Id:             user_id,
		Email:          github_email,
		Phone:          "",
		Phone_country:  "",
		Email_verify:   true,
		Phone_verify:   false,
		Avatar:         avatar,
		Username:       user_id,
		Github_connect: github_id,
		Google_connect: "",
		Email_2fa:      true,
		Phone_2fa:      false,
		Totp_2fa:       false,
		Create_at:      now_time,
		Update_at:      now_time,
	}
	fmt.Sprintln(err)
	return insert_data, err
}
