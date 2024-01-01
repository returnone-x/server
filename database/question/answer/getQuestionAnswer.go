package questionAnswerDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetQuestionAnswer(id string) ([]questionModal.QuestionAnswerModal, error) {

	sqlString := `SELECT qc.id, qc.question_id, qc.user_id, qc.content, qc.create_at, qc.update_at, u.avatar, COUNT(qavu) AS count_up_vore, COUNT(qavd) AS count_down_vore
	FROM question_answers qc
	JOIN users u ON qc.user_id = u.id
	LEFT JOIN question_answer_votes qavu ON $1 = qavu.answer_id AND qavu.vote = 1
	LEFT JOIN question_answer_votes qavd ON $1 = qavd.answer_id AND qavd.vote = 2
	WHERE qc.question_id = $1
	GROUP BY 
	qc.id, 
	qc.question_id, 
	qc.user_id, 
	qc.content, 
	qc.create_at, 
	qc.update_at, 
	u.avatar;
	`
	//	LEFT JOIN question_answer_votes uqav ON qc.user_id = uqav.voter_id AND uqav.answer_id = $1 uqav.vote
	rows, err := db.DB.Query(sqlString, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var question_answers []questionModal.QuestionAnswerModal

	for rows.Next() {
		var answer questionModal.QuestionAnswerModal

		err := rows.Scan(
			&answer.Id,
			&answer.Question_id,
			&answer.User_id,
			&answer.Content,
			&answer.Create_at,
			&answer.Update_at,
			&answer.Avatar,
			&answer.Up_vote,
			&answer.Down_vote,
		)
		if err != nil {
			return nil, err
		}

		question_answers = append(question_answers, answer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return question_answers, err
}
