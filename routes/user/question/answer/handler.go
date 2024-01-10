package questionAnswer

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	questionAnswerDatabase "github.com/returnone-x/server/database/question/answer"
	utils "github.com/returnone-x/server/utils"
)

func NewAnswer(c *fiber.Ctx) error {
	params := c.AllParams()

	// if user not send parms data
	if params["question_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	var data map[string]string
	// get data from body
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if len(data["content"]) > 100000 {
		return c.Status(400).JSON(utils.RequestValueValid("answer content"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	result, err := questionAnswerDatabase.NewQuestionAnswer(user_id, data["content"], params["question_id"])

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When create data got some error", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful post a new question answer",
		"data":    result,
	})
}

func DeleteAnswer(c *fiber.Ctx) error {
	params := c.AllParams()

	// if user not send parms data
	if params["answer_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	result, err := questionAnswerDatabase.DeleteQuestionAnswer(params["answer_id"], user_id)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When delete data got some error", err))
	}
	result_affected, _ := result.RowsAffected()
	if result_affected == 0 {
		return c.Status(400).JSON(utils.ErrorMessage("Can't find this answer or you dont have permission", err))
	}
	return c.SendStatus(204)
}

func UpVote(c *fiber.Ctx) error {
	// vote 1 = up vote
	// use function from function.go
	return Vote(c, 1)
}

func DownVote(c *fiber.Ctx) error {
	// vote 2 = donw vote
	// use function from function.go
	return Vote(c, 2)
}

func DeleteVote(c *fiber.Ctx) error {
	// get params
	params := c.AllParams()

	// if user not send parms data
	if params["answer_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	// get question data from database for check does this question exist
	answer_count, err := questionAnswerDatabase.GetQuestionAnswerByAnswerId(params["answer_id"])

	// handle error
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(utils.ErrorMessage("Can't find this answer", err))
		} else {
			return c.Status(500).JSON(utils.ErrorMessage("Error when get answer from database", err))
		}
	}

	if answer_count == 0 {
		return c.Status(400).JSON(utils.ErrorMessage("Can't find this answer", err))
	}

	// get user_id from accessToken cookie
	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	// get user vote data(this question)
	delete_result, err := questionAnswerDatabase.DeleteQuestionAnswerVote(params["answer_id"], user_id)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error delete data", err))
	}

	// check does it really update
	row_affected, _ := delete_result.RowsAffected()

	if row_affected == 0 {
		return c.Status(400).JSON(utils.ErrorMessage("Can't find this vote", err))
	} else {
		return c.SendStatus(204)
	}

}
