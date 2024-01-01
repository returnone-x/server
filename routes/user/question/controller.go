package question

import (
	"github.com/gofiber/fiber/v2"
	"github.com/returnone-x/server/routes/user/question/answer"
	questionComment "github.com/returnone-x/server/routes/user/question/comment"
)

func Setup(app fiber.Router) {
	
	// question controller
	app.Post("/new" , NewPost)
	app.Post("/upvote/:question_id", UpVote)
	app.Post("/downvote/:question_id", DownVote)
	app.Delete("/deletevote/:question_id", DeleteVote)
	
	comment_group := app.Group("/comment")
	questionComment.Setup(comment_group)

	answer_group := app.Group("/answer")
	questionAnswer.Setup(answer_group)
}