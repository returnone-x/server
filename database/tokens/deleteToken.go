package tokenDatabase

import (
	"database/sql"
	"fmt"
	db "returnone/config"
)

func DeleteToken(id string)  (sql.Result, error) {
	
	sqlString := `
	DELETE FROM tokens WHERE id = $1
	`
	reslut, err := db.DB.Exec(sqlString, id)
	fmt.Println(err)
	return reslut, err
}