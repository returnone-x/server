package questionAnswer

import (
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
