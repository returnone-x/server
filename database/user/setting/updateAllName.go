package userSettingDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateALlName(id string, display_name string, username string) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE users 
	SET 
		display_name = $2,
		username = $3,
		update_at = $4
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, display_name, username, now_time)

	return result, err
}
