package userSettingDatabase

import (
	"database/sql"

	db "github.com/returnone-x/server/config"
)

func UpdatePublicEmail(id string, bio string) (sql.Result, error) {

	sqlString := `
	UPDATE user_profile
	SET 
		public_email = $2
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, bio)

	return result, err
}
