package questionChat

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	questionChatDatabase "github.com/returnone-x/server/database/question/chat"
	utils "github.com/returnone-x/server/utils"
)

// Add more data to this type if needed
type client struct {
	question_id string
	isClosing   bool
	mu          sync.Mutex
}

type ReceivedMessage struct {
	Question_id string   `json:"question_id"`
	Method      string   `json:"method"`
	Message_id  string   `json:"message_id"`
	Content     string   `json:"content"`
	Reply       string   `json:"reply"`
	Image       []string `json:"image"`
}

var clients = make(map[*websocket.Conn]*client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var broadcast = make(chan ReceivedMessage)
var unregister = make(chan *websocket.Conn)

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = &client{question_id: connection.Params("question_id")}

		case message := <-broadcast:
			fmt.Println(message)
			for connection, c := range clients {
				fmt.Println(message)
				if c.question_id == message.Question_id {
					go func(connection *websocket.Conn, c *client) {
						c.mu.Lock()
						defer c.mu.Unlock()
						if c.isClosing {
							return
						}
						message_json, _ := json.Marshal(message)

						if err := connection.WriteMessage(websocket.TextMessage, []byte((string(message_json)))); err != nil {
							c.isClosing = true
							connection.WriteMessage(websocket.CloseMessage, []byte{})
							connection.Close()
							unregister <- connection
						}
					}(connection, c)
				}

			}

		case connection := <-unregister:
			delete(clients, connection)
		}
	}
}

func QuestionsChat(c *websocket.Conn) {
	defer func() {
		unregister <- c
		c.Close()
	}()

	question_id := c.Params("question_id")
	var userId string
	if c.Locals("access_token_context") != nil {
		token := c.Locals("access_token_context").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userId = claims["user_id"].(string)
	}

	register <- c

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}
		if messageType == websocket.TextMessage && userId != "" {
			var received_message ReceivedMessage

			err := json.Unmarshal([]byte(string(message)), &received_message)
			if err != nil {
				panic(err)
			}

			broadcast <- ReceivedMessage{
				Method:     received_message.Method,
				Message_id: received_message.Message_id,
				Content:    received_message.Content,
				Image:      received_message.Image, Question_id: question_id,
			}
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}

}

type NewMessageRequestBody struct {
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
	var data NewMessageRequestBody
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
	broadcast <- ReceivedMessage{
		Method:      "new",
		Message_id:  result_data.Id,
		Content:     result_data.Content,
		Image:       result_data.Image,
		Reply:       result_data.Reply,
		Question_id: question_id,
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "successful",
		"message": "successful send new message",
		"data":    result_data,
	})
}