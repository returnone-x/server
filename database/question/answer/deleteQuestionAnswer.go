package questionAnswerDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func DeleteQuestionAnswer(answer_id string, questioner_id string)  (sql.Result, error) {
	
	sqlString := `
	DELETE FROM question_answers WHERE id = $1 AND user_id = $2;
	`
	reslut, err := db.DB.Exec(sqlString, answer_id, questioner_id)
	return reslut, err
}