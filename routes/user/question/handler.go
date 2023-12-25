package question

import (
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	questionDatabase "github.com/returnone-x/server/database/question"
	utils "github.com/returnone-x/server/utils"
)

type RequestBody struct {
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Tags_name    []string `json:"tags_name"`
	Tags_version []string `json:"tags_version"`
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

	if len(data.Tags_name) == 0 || len(data.Tags_name) > 5 {
		return c.Status(400).JSON(utils.RequestValueValid("question tags"))
	}

	// check if tags name repeat
	seen := make(map[string]struct{})

	for _, tag := range data.Tags_name {
		if _, exists := seen[tag]; exists {
			return c.Status(400).JSON(utils.RequestValueValid("question tags"))
		}
		seen[tag] = struct{}{}
	}

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	result, err := questionDatabase.NewQuestion(user_id, data.Title, data.Content, data.Tags_name, data.Tags_version)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("When create data got some error", err))
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successful post a new question",
		"data":    result,
	})
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
