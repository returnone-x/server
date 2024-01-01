package questionModal

type QuestionVoteModal struct {
	Question_id string `json:"question_id"`
	Voter_id    string `json:"voter_id"`
	Vote        int `json:"vote"`
}

type AnswerVoteModal struct {
	Answer_id string `json:"question_id"`
	Voter_id    string `json:"voter_id"`
	Vote        int `json:"vote"`
}