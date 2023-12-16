package tokenDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
	"time"
)

func UpdateToken(id string, used_time int) (sql.Result, error) {
	// 取得當前時間
	currentTime := time.Now()

	sqlString := `
	UPDATE tokens 
	SET 
		used_time = $2,
		update_at = $3
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, used_time, currentTime)

	return result, err
}