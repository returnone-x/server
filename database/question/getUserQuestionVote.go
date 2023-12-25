package questionDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetUserQuestionVote(question_id string, voter_id string) (question_vote_data questionModal.QuestionVoteModal, err error) {
	
	sqlString := `SELECT * FROM question_votes WHERE question_id = $1 AND voter_id = $2;`
	err = db.DB.QueryRow(sqlString, question_id, voter_id).Scan(
		&question_vote_data.Question_id,
		&question_vote_data.Voter_id,
		&question_vote_data.Vote,
	)

	return question_vote_data, err
}