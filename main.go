package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id string

	answer string

	socket *websocket.Conn

	send chan []byte

	isMyTurn bool
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Test      string `json:"test,omitempty"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:

			manager.clients[conn] = true

			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected. "})
			manager.send(jsonMessage, conn)

			if len(manager.clients)+1 > 2 {
				for client := range manager.clients {
					if client == conn {
						client.isMyTurn = false
					} else {
						client.isMyTurn = true
					}
				}
			}

		case conn := <-manager.unregister:

			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected. "})
				manager.send(jsonMessage, conn)
			}

		case message := <-manager.broadcast:

			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {

		if conn != ignore {
			conn.send <- message
		}
	}
}
func (manager *ClientManager) sendMeSelf(message []byte, me *Client) {
	for conn := range manager.clients {

		if conn == me {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			_ = c.socket.Close()
			break
		}

		messageToStr := string(message)
		strSplit := strings.Split(messageToStr, "")
		isValid := true
		for _, str := range strSplit {
			_, err := strconv.Atoi(str)

			if err != nil {
				isValid = false
			}
		}

		if c.isMyTurn && isValid {
			for client := range manager.clients {
				if client == c {
					client.isMyTurn = false

					res := gameResponse(string(message), c)
					jsonMessage, _ := json.Marshal(&Message{Content: "/A Your guess number is: " + string(message)})
					manager.sendMeSelf(jsonMessage, c)
					jsonMessage, _ = json.Marshal(&Message{Content: "/A Response: " + res})
					manager.sendMeSelf(jsonMessage, c)

				} else {
					client.isMyTurn = true

					jsonMessage, _ := json.Marshal(&Message{Content: "/A User guess your number is: " + string(message)})
					manager.send(jsonMessage, c)
					jsonMessage, _ = json.Marshal(&Message{Content: "/A It's your turn. "})
					manager.send(jsonMessage, c)
				}
			}
		} else if c.isMyTurn {
			jsonMessage, _ := json.Marshal(&Message{Content: "/A Your guess number is not valid: " + string(message)})
			manager.sendMeSelf(jsonMessage, c)
		} else {
			jsonMessage, _ := json.Marshal(&Message{Content: "/A Its not your turn "})
			manager.sendMeSelf(jsonMessage, c)
		}

	}
}

func (c *Client) write() {
	defer func() {
		_ = c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.socket.WriteMessage(websocket.TextMessage, message)

		}
	}
}

func main() {
	r := gin.Default()

	fmt.Println("Starting application...")

	go manager.start()

	r.GET("/ws/:number", wsHandler)

	fmt.Println("chat server start.....")

	r.Run(":8448")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(c *gin.Context) {
	number, _ := c.Params.Get("number")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := &Client{id: uuid.Must(uuid.NewV4(), nil).String(), socket: conn, send: make(chan []byte), answer: number, isMyTurn: false}

	if len(manager.clients)+1 > 2 {
		_ = client.socket.WriteMessage(websocket.TextMessage, []byte("Room is full"))
		client.socket.Close()
		return
	}

	manager.register <- client

	go client.read()

	go client.write()
}

func gameResponse(number string, c *Client) string {
	var answer map[string]string = make(map[string]string)
	var a, b int
	var returnStr string
	for client := range manager.clients {
		if client != c {
			strAnswerList := strings.Split(client.answer, "")
			for index, item := range strAnswerList {
				answer[string(index)] = item
			}
		}
	}
	strNumberList := strings.Split(number, "")

	for index, item := range strNumberList {
		if answer[string(index)] == item {
			a += 1
		} else {
			for _, value := range answer {
				if value == item {
					b += 1
				}
			}
		}
	}
	returnStr = fmt.Sprintf("%dA %dB", a, b)

	return returnStr
}
