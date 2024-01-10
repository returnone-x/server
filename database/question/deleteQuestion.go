package questionDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func DeleteQuestion(question_id string, questioner_id string)  (sql.Result, error) {
	
	sqlString := `
	DELETE FROM questions WHERE id = $1 AND questioner_id = $2;
	`
	reslut, err := db.DB.Exec(sqlString, question_id, questioner_id)
	return reslut, err
}