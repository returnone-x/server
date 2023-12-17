package questionModal

import "time"

type QuestionModal struct {
	Id            string    `json:"id"`
	Questioner_id string    `json:"Questioner_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Tags          []TagsInfo  `json:"tags"`
	Views         int       `json:"views"`
	Create_at     time.Time `json:"create_at"`
	Update_at     time.Time `json:"update_at"`
}

type TagsInfo struct {
	Tag     string `json:"tags"`
	Version string `json:"version"`
}
