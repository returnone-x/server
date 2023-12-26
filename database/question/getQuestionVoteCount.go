package questionDatabase

import (
	db "github.com/returnone-x/server/config"
)


func GetQuestionVoteCount(question_id string, vote_value int) (int, error) {
	var voteCount int

	sqlString := `SELECT count(*) FROM question_votes WHERE question_id = $1 AND vote = $2;`
	err := db.DB.QueryRow(sqlString, question_id, vote_value).Scan(&voteCount)

	return voteCount, err
}