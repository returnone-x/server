package userDatabase

import (
	"database/sql"

	db "github.com/returnone-x/server/config"
)

func Rename(id string, new_name string) (sql.Result, error) {

	sqlString := `
	UPDATE users
	SET username = $2
	WHERE id = $1;
	`
	reslut, err := db.DB.Exec(sqlString, id, new_name)
	return reslut, err
}
