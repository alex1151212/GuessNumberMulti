package services

import (
	"gin-practice/models"
	"gin-practice/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var gameServer = models.GameServer{
	Game: &models.Game{
		Players:     make([]*models.Player, 0),
		CurrentTurn: nil,
	},
	Broadcast:  make(chan []byte),
	Register:   make(chan *models.Player),
	Unregister: make(chan *models.Player),
	Players:    make(map[*models.Player]bool),
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
			if len(gameServer.Game.Players) < 2 {
				gameServer.Game.Players = append(gameServer.Game.Players, conn)
			}

			if len(gameServer.Game.Players) == 2 {
				gameServer.Game.Init()
			}

		case conn := <-gameServer.Unregister:

			if _, ok := gameServer.Players[conn]; ok {
				close(conn.Send)
				delete(gameServer.Players, conn)
				gameServer.SendGamePlayers(utils.FormatToJson("socket has disconnected. "), conn)
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

func GameHandler(c *gin.Context) {
	number, _ := c.Params.Get("number")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	player := &models.Player{Id: uuid.Must(uuid.NewV4(), nil).String(), Socket: conn, Send: make(chan []byte), Answer: number}

	gameServer.Register <- player

	go player.Read(&gameServer)

	go player.Write(&gameServer)
}
