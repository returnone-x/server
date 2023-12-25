package question

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app fiber.Router) {
	app.Post("/new" , NewPost)
	app.Post("/upvote/:question_id", UpVote)
	app.Post("/downvote/:question_id", DownVote)
}