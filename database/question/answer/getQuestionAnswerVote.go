package questionAnswerDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetUserQuestionAnswerVote(answer_id string, voter_id string) (question_vote_data questionModal.AnswerVoteModal, err error) {
	
	sqlString := `SELECT * FROM question_answer_votes WHERE answer_id = $1 AND voter_id = $2;`
	err = db.DB.QueryRow(sqlString, answer_id, voter_id).Scan(
		&question_vote_data.Answer_id,
		&question_vote_data.Voter_id,
		&question_vote_data.Vote,
	)

	return question_vote_data, err
}