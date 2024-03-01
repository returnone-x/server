package userSettingDatabase

import (
	"database/sql"

	db "github.com/returnone-x/server/config"
)

func UpdateWebsite(id string, website string) (sql.Result, error) {
	sqlString := `
	UPDATE user_profile
	SET 
		website = $2
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, website)

	return result, err
}
