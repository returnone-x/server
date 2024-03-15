package questionChat

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pusher/pusher-http-go/v5"
	questionChatDatabase "github.com/returnone-x/server/database/question/chat"
	utils "github.com/returnone-x/server/utils"
)

// Add more data to this type if needed
type client struct {
	question_id string
	isClosing  bool
	mu         sync.Mutex
}

type Message struct {
	question_id string
	message    ReceivedMessage
}

type ReceivedMessage struct {
	Method     string   `json:"method"`
	Message_id string   `json:"message_id"`
	Message    string   `json:"message"`
	Image      []string `json:"image"`
}

var clients = make(map[*websocket.Conn]*client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var broadcast = make(chan Message)
var unregister = make(chan *websocket.Conn)

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = &client{question_id: connection.Params("questionId")}

		case message := <-broadcast:
			for connection, c := range clients {
				if c.question_id == message.question_id {
					go func(connection *websocket.Conn, c *client) {
						c.mu.Lock()
						defer c.mu.Unlock()
						if c.isClosing {
							return
						}
						if err := connection.WriteMessage(websocket.TextMessage, []byte(message.message)); err != nil {
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

	question_id := c.Params("questionId")
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
			fmt.Println("ctx:", string(received_message.Method))

			broadcast <- Message{message: received_message, question_id: question_id}
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}

}

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
