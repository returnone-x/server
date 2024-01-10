package questionDatabase

import (
	"database/sql"
	db "github.com/returnone-x/server/config"
)

func AddQuestionViews(question_id string) (sql.Result, error) {
	sqlString := `
	UPDATE questions 
	SET
	views = views + 1
	WHERE id = $1
	`
	result, err := db.DB.Exec(sqlString, question_id)

	return result, err
}
