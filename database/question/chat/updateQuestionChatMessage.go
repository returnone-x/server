package questionChatDatabase

import (
	"time"

	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func UpdateQuestionVote(message_id string, question_id, reply string, user_id string, content string, image []string) (questionModal.QuestionChat, error) {
	now_time := time.Now().UTC()
	sqlString := `
	UPDATE question_chat 
	SET 
		reply = $4, 
		content = $5, 
		image = $6, 
		update_at = $7
	WHERE 
		author = $1 AND id = $2 AND question_id = $3
	RETURNING *;
	`
	_, err := db.DB.Exec(sqlString, user_id, message_id, question_id, reply, content, pq.Array(image), now_time)

	insert_data := questionModal.QuestionChat{
		Id: message_id,
		Question_id: question_id,
		Content:     content,
		Image:       image,
		Reply:       reply,
		Create_at:   now_time,
		Update_at:   now_time,
	}
	return insert_data, err
}
