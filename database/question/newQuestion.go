package questionDatabase

import (
	"time"
	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
	utils "github.com/returnone-x/server/utils"
)

func NewQuestion(user_id string, title string, content string, tags_name []string, tags_version []string) (questionModal.QuestionModal, error) {
	question_id := utils.GenerateQuestionId()
	now_time := time.Now()
	sqlString := `
	INSERT INTO questions 
	(id, questioner_id, title, content, tags_name, tags_version, views, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.DB.Exec(sqlString, question_id, user_id, title, content, pq.Array(tags_name), pq.Array(tags_version), 0, now_time, time.Now())
	insert_data := questionModal.QuestionModal{
		Id:            question_id,
		Questioner_id: user_id,
		Title:         title,
		Content:       content,
		Tags_name:     tags_name,
		Tags_version:  tags_version,
		Views:         0,
		Create_at:     now_time,
		Update_at:     now_time,
	}
	return insert_data, err
}
