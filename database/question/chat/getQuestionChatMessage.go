package questionChatDatabase

import (
    "github.com/lib/pq"
    db "github.com/returnone-x/server/config"
    questionModal "github.com/returnone-x/server/models/question"
)

func GetChatQuestionhatMessage(question_id string, page string) ([]questionModal.QuestionChat, error) {
    sqlString := `
        SELECT 
            q.id, 
            q.question_id, 
            q.author,
            q.content,
            q.image,
            q.reply, 
            q.create_at, 
            q.update_at
        FROM 
            question_chat q
        WHERE 
            q.question_id = $1
        LIMIT 50
        OFFSET ($2 - 1) * 50;
    `

    rows, err := db.DB.Query(sqlString, question_id, page)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []questionModal.QuestionChat
    for rows.Next() {
        var message questionModal.QuestionChat
        err := rows.Scan(
            &message.Id,
            &message.Question_id,
            &message.Author,
            &message.Content,
            (*pq.StringArray)(&message.Image),
            &message.Reply,
            &message.Create_at,
            &message.Update_at,
        )
        if err != nil {
            return nil, err
        }

        messages = append(messages, message)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}