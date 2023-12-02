package userModles

import "time"

type UserAccount struct {
	Id             string    `json:"id"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Phone_country  string    `json:"phone_country"`
	Password       string    `json:"password"`
	Email_verify   bool      `json:"email_verify"`
	Phone_verify   bool      `json:"phone_verify"`
	Avatar         string    `json:"avatar"`
	User_name      string    `json:"user_name"`
	Github_connect string    `json:"github_connect"`
	Google_connect string    `json:"google_connect"`
	Create_at      time.Time `json:"create_at"`
	Update_at      time.Time `json:"update_at"`
}

// CREATE TABLE IF NOT EXISTS users (
//     id VARCHAR(255) PRIMARY KEY,
//     email VARCHAR(100) NOT NULL,
//     phone VARCHAR(50),
//     phone_country VARCHAR(25),
//     password VARCHAR(255) NOT NULL,
//     email_verify BOOLEAN DEFAULT FALSE,
//     phone_verify BOOLEAN DEFAULT FALSE,
//     avatar VARCHAR(255),
//     user_name VARCHAR(30) NOT NULL,
//     github_connect VARCHAR(100),
//     google_connect VARCHAR(100),
//     create_at TIMESTAMP,
//     update_at TIMESTAMP
// );
