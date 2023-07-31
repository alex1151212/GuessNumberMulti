package services

import (
	"gin-practice/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var gameServer = models.GameServer{
	Game:       make(map[string]*models.Game),
	Broadcast:  make(chan []byte),
	Register:   make(chan *models.Player),
	Unregister: make(chan *models.Player),
	Players:    make(map[string]*models.Player),
}

// &models.Game{
// 	Players:     make([]*models.Player, 0),
// 	CurrentTurn: nil,
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GameServerStart() {
	gameServer.Init()
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

func CreateGame(c *gin.Context) {
	gameId := uuid.Must(uuid.NewV4(), nil).String()
	gameServer.Game[gameId] = &models.Game{
		Id:          uint(len(gameServer.Game)) + 1,
		Players:     make([]*models.Player, 0),
		CurrentTurn: nil,
	}

	//
	c.JSON(http.StatusOK, gin.H{
		"GamesList": gameServer.Game,
	})
}

func GetOnlinePlayers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"PlayerList": gameServer.Players,
	})
}
func GetGames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"GameList": gameServer.Game,
	})
}

func DeleteGame(c *gin.Context) {
	gameId, _ := c.Params.Get("gameId")

	delete(gameServer.Game, gameId)
	//
	c.JSON(http.StatusOK, gin.H{
		"GamesList": gameServer.Game,
	})
}

func JoinGame(c *gin.Context) {
	gameId, _ := c.Params.Get("gameId")
	playerId, _ := c.Params.Get("playerId")

	gameServer.Game[gameId].Players = append(gameServer.Game[gameId].Players, gameServer.Players[playerId])
	gameServer.Players[playerId].GameId = gameId

	//
	c.JSON(http.StatusOK, gin.H{
		"status": "Successfully joined",
	})
}
