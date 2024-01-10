package questionAnswerDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
	utils "github.com/returnone-x/server/utils"
	"time"
)

func NewQuestionAnswer(user_id string, content string, question_id string) (questionModal.QuestionAnswerModal, error) {
	answer_id := utils.GenerateQuestionId()
	now_time := time.Now().UTC()
	sqlString := `
	INSERT INTO question_answers
	(id, question_id, content, user_id, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6)
	`
	_, err := db.DB.Exec(sqlString, answer_id, question_id, content, user_id, now_time, now_time)
	
	insert_data := questionModal.QuestionAnswerModal{
		Id:          answer_id,
		Question_id: question_id,
		User_id:     user_id,
		Content:     content,
		Create_at:   now_time,
		Update_at:   now_time,
	}
	return insert_data, err
}
