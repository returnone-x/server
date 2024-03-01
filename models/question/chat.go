package questionModal

import "time"

type QuestionChat struct {
	Id          string    `json:"id"`
	Question_id string    `json:"question_id"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	Image       []string  `json:"image"`
	Reply       string    `json:"reply"`
	Create_at   time.Time `json:"create_at"`
	Update_at   time.Time `json:"update_at"`
}
