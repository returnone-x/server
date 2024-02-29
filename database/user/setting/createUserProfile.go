package userSettingDatabase

import (
	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
)

func CreateUserProfile(user_id string) error {
	sqlString := `
	INSERT INTO user_profile
	(id, bio, public_email, pronouns, website, related_links) 
	VALUES 
	($1, $2, $3, $4, $5, $6)
	`
	_, err := db.DB.Exec(sqlString, user_id, "", "", "", "", pq.Array([]string{}))

	return err
}
