package userSettingDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateTotp(id string, totp string) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE users 
	SET 
		totp = $2,
		update_at = $3
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, totp, now_time)

	return result, err
}
