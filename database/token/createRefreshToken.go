package tokenDatabase

import (
	"returnone/config"
)

func CreateRefreshToken(
	token string,
) (int64, error) {

	sqlString := `
	INSERT INTO token 
	(token, used) 
	VALUES 
	($1, $2)
	`
	result, err := config.DB.Exec(
		sqlString,
		token, false)
	
	affected, _ := result.RowsAffected()
	
	return affected, err
}
