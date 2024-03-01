package questionComment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	questionCommentDatabase "github.com/returnone-x/server/database/question/comment"
	utils "github.com/returnone-x/server/utils"
)

func NewComment(c *fiber.Ctx) error {
	params := c.AllParams()

	// if user not send parms data
	if params["question_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if len(data["content"]) > 3000 || len(data["content"]) == 3000 {
		return c.Status(400).JSON(utils.RequestValueValid("question content"))
	}

	if len(data["reply"]) > 50 {
		return c.Status(400).JSON(utils.RequestValueValid("question tags"))
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	result, err := questionCommentDatabase.NewQuestionComment(params["question_id"], user_id, data["content"], data["reply"])

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error create data", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "successful create new comment data",
		"data":    result,
	})
}
