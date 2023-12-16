package tokenDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func CreateToken(id string,user_agent string, ip string)  (sql.Result, error) {
	
	sqlString := `
	INSERT INTO tokens 
	(id, used_time, user_agent, ip, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6)
	`
	reslut, err := db.DB.Exec(sqlString, id, 1, user_agent, ip, time.Now(), time.Now())

	return reslut, err
}