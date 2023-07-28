package services

import (
	"fmt"
	"gin-practice/models"
	"gin-practice/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type GameServer struct {
	Players    map[*models.Player]bool
	Game       *models.Game
	Register   chan *models.Player
	Unregister chan *models.Player
	Broadcast  chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Init() {
	for {
		select {
		case conn := <-gameServer.Register:

			gameServer.Players[conn] = true

			gameServer.SendPlayers(utils.FormatToJson("A new socket has connected. "), conn)

			// 暫時邏輯 若gameserver 人數大於2 則創建遊戲

			if len(gameServer.Game.Players) == 2 {
				gameServer.Game.Init()
				// gameServer.
			} else {
				gameServer.Game.Players = append(gameServer.Game.Players, conn)
			}

		case conn := <-gameServer.Unregister:

			if _, ok := gameServer.Players[conn]; ok {
				close(conn.Send)
				delete(gameServer.Players, conn)
				// jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected. "})
				// game.Send(jsonMessage, conn)
			}

		case message := <-gameServer.Broadcast:

			for conn := range gameServer.Players {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(gameServer.Players, conn)
				}
			}
		}
	}
}

// 發送給指定遊戲內的所有玩家
func (gameServer *GameServer) SendGamePlayers(message []byte, ignore *models.Player) {
	for _, conn := range gameServer.Game.Players {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 發送給指定遊戲內的指定玩家
func (gameServer *GameServer) SendGamePlayer(message []byte, player *models.Player) {
	for _, conn := range gameServer.Game.Players {
		if conn == player {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的所有玩家
func (gameServer *GameServer) SendPlayers(message []byte, ignore *models.Player) {
	for conn := range gameServer.Players {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的指定玩家
func (gameServer *GameServer) SendPlayer(message []byte, player *models.Player) {
	for conn := range gameServer.Players {
		if conn == player {
			conn.Send <- message
		}
	}
}

var gameServer = GameServer{
	Game: &models.Game{
		Players:     make([]*models.Player, 0),
		CurrentTurn: nil,
	},
	Broadcast:  make(chan []byte),
	Register:   make(chan *models.Player),
	Unregister: make(chan *models.Player),
	Players:    make(map[*models.Player]bool),
}

func GameHandler(c *gin.Context) {
	number, _ := c.Params.Get("number")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	player := &models.Player{Id: uuid.Must(uuid.NewV4(), nil).String(), Socket: conn, Send: make(chan []byte), Answer: number}

	gameServer.Register <- player

	go player.Read()

	go player.Write()
}

func gameResponse(number string, c *models.Player) string {
	var answer map[string]string = make(map[string]string)
	var a, b int
	var returnStr string
	for _, client := range gameServer.Game.Players {
		if client != c {
			strAnswerList := strings.Split(client.Answer, "")
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
