package questionAnswerDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func NewQuestionAnswerVote(answer_id string, voter_id string, vote int) (questionModal.AnswerVoteModal, error) {
	sqlString := `
	INSERT INTO question_answer_votes 
	(answer_id, voter_id, vote) 
	VALUES 
	($1, $2, $3)
	`
	_, err := db.DB.Exec(sqlString, answer_id, voter_id, vote)
	insert_data := questionModal.AnswerVoteModal{
		Answer_id: answer_id,
		Voter_id:    voter_id,
		Vote:        vote,
	}
	return insert_data, err
}
