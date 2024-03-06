package questionDatabase

import (
	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetQuestionData(id string) (question_data questionModal.ReturnSourceResult, err error) {
	sqlString := `
	SELECT 
		q.id, 
		q.questioner_id, 
		q.title,
		q.tags_name,
		q.tags_version,
		q.content, 
		q.views, 
		q.create_at,
		q.update_at,
		u.username,
		u.avatar
	FROM 
		questions q
	JOIN 
		users u ON q.questioner_id = u.id
	WHERE 
		q.id = $1
	GROUP BY 
		q.id, 
		q.questioner_id, 
		q.title,
		q.tags_name,
		q.tags_version,
		q.content,
		q.views, 
		q.create_at,
		q.update_at,
		u.avatar,
		u.username;
	`
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
		&question_data.Questioner_name,
		&question_data.Questioner_avatar,
	)

	return question_data, err
}
