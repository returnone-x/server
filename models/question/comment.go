package questionModal

import "time"

type QuestionCommentModal struct {
	Id           string    `json:"id"`
	Question_id  string    `json:"question_id"`
	Commenter_id string    `json:"commenter_id"`
	Content      string    `json:"content"`
	Reply        string    `json:"reply"`
	Update_at    time.Time `json:"update_at"`
	Create_at    time.Time `json:"create_at"`
	User_id      string    `json:"user_id"`
	Avatar       string    `json:"avatar"`
}
