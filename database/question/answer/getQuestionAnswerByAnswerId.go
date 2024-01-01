package questionAnswerDatabase

import (
	db "github.com/returnone-x/server/config"
)

func GetQuestionAnswerByAnswerId(answer_id string) (int, error) {
	var answer_count int

	sqlString := `SELECT count(*) FROM question_answers qc WHERE id = $1;`

	err := db.DB.QueryRow(sqlString, answer_id).Scan(&answer_count)

	return answer_count, err
}
