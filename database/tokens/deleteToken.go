package tokenDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func DeleteToken(id string)  (sql.Result, error) {
	
	sqlString := `
	DELETE FROM tokens WHERE id = $1
	`
	reslut, err := db.DB.Exec(sqlString, id)
	return reslut, err
}