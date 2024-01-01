package questionAnswerDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func UpdateQuestionAnswerVote(answer_id string, voter_id string, vote int) (sql.Result, error) {
	sqlString := `
	UPDATE question_answer_votes 
	SET 
	vote = $3
	WHERE answer_id = $1 AND voter_id = $2
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, answer_id, voter_id, vote)

	return result, err
}
