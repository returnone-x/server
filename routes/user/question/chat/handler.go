package questionChat

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pusher/pusher-http-go/v5"
	questionChatDatabase "github.com/returnone-x/server/database/question/chat"
	utils "github.com/returnone-x/server/utils"
)

type RequestBody struct {
	Reply   string   `json:"reply"`
	Content string   `json:"content"`
	Image   []string `json:"image"`
}


func NewMessage(c *fiber.Ctx) error {
	params := c.AllParams()

	// if user not send parms data
	if params["question_id"] == "" {
		return c.Status(400).JSON(utils.InvalidRequest())
	}
	question_id := params["question_id"]
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

	token := c.Locals("access_token_context").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	// get user_id from accessToken cookie
	user_id := claims["user_id"].(string)

	if len(data.Image) > 5 {
		return c.Status(400).JSON(utils.RequestValueValid("image"))
	}

	if len(data.Content) > 3000 || len(data.Content) < 1 {
		return c.Status(400).JSON(utils.RequestValueValid("content"))
	}

	if len(data.Image) < 1 {
		data.Image = []string{}
	}

	result_data, err := questionChatDatabase.NewQuestionChatMessage(question_id, data.Reply, user_id, data.Content, data.Image)

	if err != nil {
		return c.Status(500).JSON(utils.ErrorMessage("Error save data into database", err))
	}

	pusherClient := pusher.Client{
		AppID:   os.Getenv("APP_ID"),
		Key:     os.Getenv("KEY"),
		Secret:  os.Getenv("SECRET"),
		Cluster: os.Getenv("CLUSTER"),
		Secure:  true,
	}
	
	err = pusherClient.Trigger(question_id, "new-message", result_data)

	if err != nil {
	  fmt.Println(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful send new message",
		"data":    result_data,
	})
}