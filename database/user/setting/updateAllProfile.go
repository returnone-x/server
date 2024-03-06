package userSettingDatabase

import (
	"database/sql"

	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
)

func UpdateAllProfile(id string, bio string, public_email string, pronouns string, related_links []string) (sql.Result, error) {

	sqlString := `
	UPDATE user_profile
	SET 
		bio = $2,
		public_email = $3,
		pronouns = $4,
		related_links = $5
	WHERE id = $1
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, id, bio, public_email, pronouns, pq.Array(related_links))

	return result, err
}
