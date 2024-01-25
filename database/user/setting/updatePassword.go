package userSettingDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateUserPassword(id string, new_password string) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE users 
	SET 
		Password = $2,
		update_at = $3
	WHERE id = $1a
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, new_password, now_time)

	return result, err
}
