package userSettingDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateDisplayName(id string, display_name string) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE users 
	SET 
		display_name = $2,
		update_at = $3
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, display_name, now_time)

	return result, err
}
