package questionCommentDatabase

import (
	"time"

	db "github.com/returnone-x/server/config"
	userDatabase "github.com/returnone-x/server/database/user"
	questionModal "github.com/returnone-x/server/models/question"
	utils "github.com/returnone-x/server/utils"
)

func NewQuestionComment(question_id string, commenter_id string, content string, reply string) (questionModal.QuestionCommentModal, error) {
	comment_id := utils.GenerateQuestionId()
	now_time := time.Now()
	sqlString := `
	INSERT INTO question_comments 
	(id, question_id, commenter_id, content, reply, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := db.DB.Exec(sqlString, comment_id, question_id, commenter_id, content, reply, now_time, time.Now())

	if err != nil{
		insert_data := questionModal.QuestionCommentModal{
			Id:            comment_id,
			Question_id: question_id,
			Content:       content,
			Reply:         reply,
			Create_at:     now_time,
			Update_at:     now_time,
		}
		return insert_data, err
	}

	result, err := userDatabase.GetUserAvatar(commenter_id)

	insert_data := questionModal.QuestionCommentModal{
		Id:            comment_id,
		Question_id: question_id,
		Content:       content,
		Reply:         reply,
		Create_at:     now_time,
		Update_at:     now_time,
		Avatar: result,
	}
	return insert_data, err
}
