package questionAnswerDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func DeleteQuestionAnswerVote(answer_id string, voter_id string)  (sql.Result, error) {
	
	sqlString := `
	DELETE FROM question_answer_votes WHERE answer_id = $1 AND voter_id = $2;
	`
	reslut, err := db.DB.Exec(sqlString, answer_id, voter_id)
	return reslut, err
}