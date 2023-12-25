package public

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	questionDatabase "github.com/returnone-x/server/database/question"
	utils "github.com/returnone-x/server/utils"
)

func GetQuestion(c *fiber.Ctx) error {

	params := c.AllParams()

	if params["id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}
	question_data, err := questionDatabase.GetQuestionData(params["id"])

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(utils.ErrorMessage("Can't find this question", err))
		} else {
			return c.Status(500).JSON(utils.ErrorMessage("Error when get question from database", err))
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful get the question data",
		"data":    question_data,
	})
}
