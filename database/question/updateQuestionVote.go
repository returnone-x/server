package questionDatabase

import (
	"database/sql"

	db "github.com/returnone-x/server/config"
)

func UpdateQuestionVote(question_id string, voter_id string, vote int) (sql.Result, error) {
	sqlString := `
	UPDATE question_votes 
	SET 
	vote = $3
	WHERE question_id = $1 AND voter_id = $2
	RETURNING *
	`
	result, err := db.DB.Exec(sqlString, question_id, voter_id, vote)

	return result, err
}
