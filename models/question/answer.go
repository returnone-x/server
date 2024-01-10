package questionModal

import "time"

type QuestionAnswerModal struct {
	Id          string    `json:"id"`
	Question_id string    `json:"question_id"`
	User_id     string    `json:"user_id"`
	Content     string    `json:"content"`
	Avatar      string    `json:"avatar"`
	User_name   string    `json:"user_name"`
	Up_vote     int       `json:"up_vote"`
	Down_vote   int       `json:"down_vote"`
	User_vote   int       `json:"user_vote"`
	Create_at   time.Time `json:"create_at"`
	Update_at   time.Time `json:"update_at"`
}
