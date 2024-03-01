package userSettingDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func UpdatePronouns(id string, pronouns string) (sql.Result, error) {

	sqlString := `
	UPDATE user_profile
	SET 
	pronouns = $2
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, pronouns)

	return result, err
}
