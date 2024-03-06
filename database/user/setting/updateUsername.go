package userSettingDatabase

import (
	"database/sql"
	"time"

	db "github.com/returnone-x/server/config"
)

func UpdateUsername(id string, username string) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE users 
	SET 
		username = $2,
		update_at = $3
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, username, now_time)

	return result, err
}
