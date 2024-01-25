package questionModal

import "time"

type QuestionModal struct {
	Id            string    `json:"id"`
	Questioner_id string    `json:"Questioner_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Tags_name     []string  `json:"tags_name"`
	Tags_version  []string  `json:"tags_version"`
	Views         int       `json:"views"`
	Create_at     time.Time `json:"create_at"`
	Update_at     time.Time `json:"update_at"`
}

type ReturnSourceResult struct {
	Id                string    `json:"id"`
	Questioner_id     string    `json:"Questioner_id"`
	Title             string    `json:"title"`
	Content           string    `json:"content"`
	Tags_name         []string  `json:"tags_name"`
	Tags_version      []string  `json:"tags_version"`
	Views             int       `json:"views"`
	Create_at         time.Time `json:"create_at"`
	Update_at         time.Time `json:"update_at"`
	Questioner_name   string    `json:"questioner_name"`
	Questioner_avatar string    `json:"questioner_avatar"`
}

type ReturnResultWithVoteAndAnswers struct {
	Id                 string    `json:"id"`
	Questioner_id      string    `json:"Questioner_id"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	Tags_name          []string  `json:"tags_name"`
	Tags_version       []string  `json:"tags_version"`
	Views              int       `json:"views"`
	Create_at          time.Time `json:"create_at"`
	Update_at          time.Time `json:"update_at"`
	Questioner_name    string    `json:"questioner_name"`
	Questioner_avatar  string    `json:"questioner_avatar"`
	Question_vote_down string    `json:"vote_down"`
	Question_vote_up   string    `json:"vote_up"`
	Question_answers   int       `json:"answers"`
}

type TagsInfo struct {
	Tag     string `json:"tags"`
	Version string `json:"version"`
}

type ReturnResult struct {
	Id                string                `json:"id"`
	Questioner_id     string                `json:"questioner_id"`
	Title             string                `json:"title"`
	Content           string                `json:"content"`
	Tags_name         []string              `json:"tags_name"`
	Tags_version      []string              `json:"tags_version"`
	Views             int                   `json:"views"`
	Answers           []QuestionAnswerModal `json:"answers"`
	Vote_count        int                   `json:"vote_count"`
	User_vote         int                   `json:"user_vote"`
	Create_at         time.Time             `json:"create_at"`
	Update_at         time.Time             `json:"update_at"`
	Questioner_name   string                `json:"questioner_name"`
	Questioner_avatar string                `json:"questioner_avatar"`
}
