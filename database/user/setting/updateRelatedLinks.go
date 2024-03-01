package userSettingDatabase

import (
	"database/sql"

	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
)

func UpdateRelateLinks(id string, related_links []string) (sql.Result, error) {

	sqlString := `
	UPDATE user_profile
	SET 
		related_links = $2
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, pq.Array(related_links))

	return result, err
}
