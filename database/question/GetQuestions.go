package questionDatabase

import (
	"github.com/lib/pq"
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetQuestions(page int) ([]questionModal.ReturnResultWithVoteAndAnswers, error) {
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
		u.avatar,
		COUNT(qvu) AS count_up_vote, 
		COUNT(qvd) AS count_down_vote,
		COUNT(qa) AS count_question_answers
	FROM 
		questions q
	LEFT JOIN
		users u ON q.questioner_id = u.id
	LEFT JOIN 
		question_votes qvu ON q.id = qvu.question_id AND qvu.vote = 1
	LEFT JOIN 
		question_votes qvd ON q.id = qvd.question_id AND qvd.vote = 2
	LEFT JOIN 
		question_answers qa ON q.id = qa.question_id
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
		u.username
	ORDER BY q.create_at DESC
	LIMIT 15
	OFFSET ($1 - 1) * 15;
	`

	rows, err := db.DB.Query(sqlString, page)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []questionModal.ReturnResultWithVoteAndAnswers

	for rows.Next() {
		var question questionModal.ReturnResultWithVoteAndAnswers

		err := rows.Scan(
			&question.Id,
			&question.Questioner_id,
			&question.Title,
			(*pq.StringArray)(&question.Tags_name),
			(*pq.StringArray)(&question.Tags_version),
			&question.Content,
			&question.Views,
			&question.Create_at,
			&question.Update_at,
			&question.Questioner_name,
			&question.Questioner_avatar,
			&question.Question_vote_up,
			&question.Question_vote_down,
			&question.Question_answers,
		)
		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, err
}

func GetQuestionsNumber() (int, error) {
	var questions_count int

	sqlString := `SELECT count(*) FROM questions;`

	err := db.DB.QueryRow(sqlString).Scan(&questions_count)

	return questions_count, err
}
