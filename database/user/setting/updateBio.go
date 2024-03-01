package userSettingDatabase

import (
	"database/sql"

	db "github.com/returnone-x/server/config"
)

func UpdateBio(id string, bio string) (sql.Result, error) {

	sqlString := `
	UPDATE user_profile
	SET 
		bio = $2
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, bio)

	return result, err
}
