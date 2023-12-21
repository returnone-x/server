package questionDatabase

import (
	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetQuestionData(id string) (question_data questionModal.QuestionModal, err error) {
	
	sqlString := `SELECT * FROM questions WHERE id = $1;`
	err = db.DB.QueryRow(sqlString, id).Scan(
		&question_data.Id,
		&question_data.Questioner_id,
		&question_data.Title,
		(*pq.StringArray)(&question_data.Tags_name),
		(*pq.StringArray)(&question_data.Tags_version),
		&question_data.Content,
		&question_data.Views,
		&question_data.Create_at,
		&question_data.Update_at,
	)

	return question_data, err
}
