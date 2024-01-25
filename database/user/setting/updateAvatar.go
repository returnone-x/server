package userSettingDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateUserAvatar(id string, avatar string) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE users 
	SET 
		avatar = $2,
		update_at = $3
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, avatar, now_time)

	return result, err
}
