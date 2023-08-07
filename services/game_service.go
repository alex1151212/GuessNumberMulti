package services

import (
	"encoding/json"
	"fmt"
	messageType "gin-practice/enum/message"
	playerStatusType "gin-practice/enum/playerStatus"
	"gin-practice/models"
	"gin-practice/utils"
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
	// gameId := uuid.Must(uuid.NewV4(), nil).String()
	gameId, _ := c.Params.Get("gameId")

	if gameServer.Game[gameId] == nil {
		gameServer.Game[gameId] = &models.Game{
			Id:          gameId,
			Players:     make([]*models.Player, 0),
			CurrentTurn: nil,
		}
	}

	type gameResponse struct {
		Id           string `json:"id"`
		PlayerAmount int    `json:"playerAmount"`
	}
	var gameList = make(map[string]*gameResponse)
	for k, v := range gameServer.Game {
		gameList[k] = &gameResponse{
			Id:           k,
			PlayerAmount: len(v.Players),
		}
	}
	gameRespData, err := json.Marshal(gameList)
	if err != nil {
		fmt.Println(gameRespData)
	}
	for _, player := range gameServer.Players {
		if player.Status == playerStatusType.INLOBBY {
			gameServer.SendPlayer(utils.RespMessage(&utils.Message{
				Type: messageType.GET_GAMES,
				Data: string(gameRespData),
			}), player)
		}
	}

	//
	c.JSON(http.StatusOK, gin.H{
		"GamesList": gameServer.Game,
	})
}

func GetOnlinePlayers(c *gin.Context) {

	var players []string = make([]string, 0)
	for key, _ := range gameServer.Players {
		players = append(players, key)
	}
	c.JSON(http.StatusOK, gin.H{
		"PlayerList": players,
	})
}
func GetGames(c *gin.Context) {
	var gameList = make([]*models.Game, 0)
	for _, v := range gameServer.Game {
		gameList = append(gameList, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"GameList": gameList,
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

	game := gameServer.Game[gameId]

	game.Players = append(game.Players, gameServer.Players[playerId])
	gameServer.Players[playerId].GameId = gameId

	if len(game.Players) == 2 {
		game.Init()
		gameServer.SendGamePlayers(gameId, utils.Resp("Game Start"), nil)
	}

	//
	c.JSON(http.StatusOK, gin.H{
		"status": "Successfully joined",
	})
}
