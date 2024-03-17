package questionChatDatabase

import (
	"time"

	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
	utils "github.com/returnone-x/server/utils"
)

func NewQuestionChatMessage(question_id string, reply string, author string, content string, image []string) (questionModal.QuestionChat, error) {
	message_id := utils.GenerateQuestionChatId()
	now_time := time.Now().UTC()
	sqlString := `
	INSERT INTO question_chat 
	(id, question_id, reply, author, content, image, create_at, update_at) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := db.DB.Exec(sqlString, message_id, question_id, reply, author, content, pq.Array(image), now_time, now_time)

	insert_data := questionModal.QuestionChat{
		Id: message_id,
		Question_id: question_id,
		Author:      author,
		Content:     content,
		Image:       image,
		Reply:       reply,
		Create_at:   now_time,
		Update_at:   now_time,
	}
	return insert_data, err
}
