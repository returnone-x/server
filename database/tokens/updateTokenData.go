package tokenDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateToken(id string, used_time int) (sql.Result, error) {
	now_time := time.Now().UTC()

	sqlString := `
	UPDATE tokens 
	SET 
		used_time = $2,
		update_at = $3
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, used_time, now_time)

	return result, err
}