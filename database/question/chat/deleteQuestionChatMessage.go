package questionChatDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func DeleteQuestionMessage(question_id string, message_id string, user_id string) (sql.Result, error) {

	sqlString := `
	DELETE FROM question_chat WHERE id = $1 AND question_id = $2 AND author = $3;
	`
	reslut, err := db.DB.Exec(sqlString, message_id, question_id, user_id)
	return reslut, err
}
