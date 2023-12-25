package questionDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func NewQuestionVote(question_id string, voter_id string, vote int) (questionModal.QuestionVoteModal, error) {
	sqlString := `
	INSERT INTO question_votes 
	(question_id, voter_id, vote) 
	VALUES 
	($1, $2, $3)
	`
	_, err := db.DB.Exec(sqlString, question_id, voter_id, vote)
	insert_data := questionModal.QuestionVoteModal{
		Question_id: question_id,
		Voter_id:    voter_id,
		Vote:        vote,
	}
	return insert_data, err
}
