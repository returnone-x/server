package questionChat

import (
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
	isClosing bool
	mu        sync.Mutex
}

var clients = make(map[*websocket.Conn]*client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = &client{}
			log.Println("connection registered")

		case message := <-broadcast:
			log.Println("message received:", message)
			// Send the message to all clients
			for connection, c := range clients {
				go func(connection *websocket.Conn, c *client) { // send to each client in parallel so we don't block on a slow client
					c.mu.Lock()
					defer c.mu.Unlock()
					if c.isClosing {
						return
					}
					if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						c.isClosing = true
						log.Println("write error:", err)

						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						unregister <- connection
					}
				}(connection, c)
			}

		case connection := <-unregister:
			// Remove the client from the hub
			delete(clients, connection)

			log.Println("connection unregistered")
		}
	}
}

func QuestionsChat(c *websocket.Conn) {
	defer func() {
		unregister <- c
		c.Close()
	}()
	// c.Locals is added to the *websocket.Conn
	question_id := c.Params("questionId")
	var user_id string
	if c.Locals("access_token_context") == nil {

	} else {
		token := c.Locals("access_token_context").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		user_id = claims["user_id"].(string)
	}

	fmt.Println(user_id)
	fmt.Println(question_id)

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
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
