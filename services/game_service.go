package services

import (
	gameStatusType "guessNumber/enum/gameStatus"
	messageType "guessNumber/enum/message"
	playerStatusType "guessNumber/enum/playerStatus"
	"guessNumber/models"
	"guessNumber/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var gameServer = models.GameServer{
	Game:       make(map[string]*models.Game),
	Broadcast:  make(chan []byte, 1),
	Register:   make(chan *models.Player),
	Unregister: make(chan *models.Player),
	Players:    make(map[string]*models.Player),
	GameNew:    make(chan string),
	GameEnd:    make(chan *models.Game),
}

var upgrader = websocket.Upgrader{
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GameServerStart() {
	gameServer.Init()
}

func GameHandler(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	player := &models.Player{Id: uuid.Must(uuid.NewV4(), nil).String(), Socket: conn, Send: make(chan []byte), Answer: nil}

	gameServer.Register <- player

	go player.Read(&gameServer)

	go player.Write(&gameServer)

}

/*
測試用 API
*/

func CreateGame(c *gin.Context) {

	gameId, _ := c.Params.Get("gameId")
	game, ok := gameServer.Game[gameId]

	if !ok {
		gameServer.Game[gameId] = &models.Game{
			Id:          &gameId,
			Players:     make(map[string]*models.Player),
			CurrentTurn: nil,
			Winner:      nil,
			Status:      gameStatusType.WAITING,
			Broadcast:   make(chan []byte, 1),
		}
		game = gameServer.Game[gameId]
	}
	game.Init(&gameServer)
	var gameList = make(map[string]*utils.GameRoomRespType)
	for k, v := range gameServer.Game {
		gameList[k] = &utils.GameRoomRespType{
			Id:           k,
			PlayerAmount: len(v.Players),
		}
	}
	for _, player := range gameServer.Players {
		if player.Status == playerStatusType.INLOBBY {
			gameServer.SendPlayer(utils.RespMessage(messageType.GET_GAMES, gameList), player)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"GamesList": gameServer.Game,
	})
}

func GetOnlinePlayers(c *gin.Context) {

	var players []string = make([]string, 0)
	for key := range gameServer.Players {
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

	c.JSON(http.StatusOK, gin.H{
		"GamesList": gameServer.Game,
	})
}

func JoinGame(c *gin.Context) {
	gameId, _ := c.Params.Get("gameId")
	playerId, _ := c.Params.Get("playerId")

	game := gameServer.Game[gameId]

	gameServer.Players[playerId].JoinGame(game)
	gameServer.Players[playerId].Game = game

	if len(game.Players) == 2 {
		game.Init(&gameServer)
		gameServer.Game[*game.Id].Broadcast <- utils.Resp("Game Start")
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Successfully joined",
	})
}
