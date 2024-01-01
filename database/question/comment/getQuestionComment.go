package questionCommentDatabase

import (
	db "github.com/returnone-x/server/config"
	questionModal "github.com/returnone-x/server/models/question"
)

func GetQuestionComments(question_id string, limit int) ([]questionModal.QuestionCommentModal, error) {
	
	sqlString := `
	SELECT qc.id, qc.question_id, qc.commenter_id, qc.content, qc.reply, qc.create_at, qc.update_at, u.avatar, u.id
	FROM question_comments qc
	JOIN users u ON qc.commenter_id = u.id
	WHERE qc.question_id = $1
	ORDER BY qc.create_at ASC
	LIMIT $2
	`

	rows, err := db.DB.Query(sqlString, question_id, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var question_comments []questionModal.QuestionCommentModal

	for rows.Next() {
		var comment questionModal.QuestionCommentModal

		err := rows.Scan(
			&comment.Id,
			&comment.Question_id,
			&comment.Commenter_id,
			&comment.Content,
			&comment.Reply,
			&comment.Create_at,
			&comment.Update_at,
			&comment.User_id,
			&comment.Avatar,
		)
		if err != nil {
			return nil, err
		}

		question_comments = append(question_comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return question_comments, err
}