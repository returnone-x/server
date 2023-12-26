package questionDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func DeleteQuestionVote(question_id string, voter_id string)  (sql.Result, error) {
	
	sqlString := `
	DELETE FROM question_votes WHERE question_id = $1 AND voter_id = $2;
	`
	reslut, err := db.DB.Exec(sqlString, question_id, voter_id)
	return reslut, err
}