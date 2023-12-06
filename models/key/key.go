package keyModles

import "time"

type TokenType struct {
	Id         string    `json:"Id"`
	Used_time   int       `json:"Used_time"`
	User_agent string    `json:"User_agent"`
	Ip         string    `json:"Ip"`
	Create_at  time.Time `json:"Create_at"`
	Update_at  time.Time `json:"Update_at"`
}
