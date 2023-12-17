package question

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	questionDatabase "github.com/returnone-x/server/database/question"
	questionModal "github.com/returnone-x/server/models/question"
	utils "github.com/returnone-x/server/utils"
)

type RequestBody struct {
	Title   string                   `json:"title"`
	Content string                   `json:"content"`
	Tags    []questionModal.TagsInfo `json:"tags"`
}

func NewPost(c *fiber.Ctx) error {
	var data RequestBody
	// get data from body
	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid post request",
			})
	}

	if len(data.Title) > 200 {
		return c.Status(400).JSON(utils.RequestValueValid("question title"))
	}

	if len(data.Content) > 100000 {
		return c.Status(400).JSON(utils.RequestValueValid("question content"))
	}

	if len(data.Tags) == 0 || len(data.Tags) > 5 {
		return c.Status(400).JSON(utils.RequestValueValid("question tags"))
	}

	fmt.Println(data.Tags)

	// check if tags name repeat
	seen := make(map[string]struct{})

	for _, tag := range data.Tags {
		if _, exists := seen[tag.Tag]; exists {
			return c.Status(400).JSON(utils.RequestValueValid("question tags"))
		}
		seen[tag.Tag] = struct{}{}
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	result, err := questionDatabase.NewQuestion(user_id, data.Title, data.Content, data.Tags)
	fmt.Println(err)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When create data got some error", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
