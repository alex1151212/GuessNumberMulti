package models

import (
	"encoding/json"
	"fmt"
	messageType "guessNumber/enum/message"
	playerStatusType "guessNumber/enum/playerStatus"
	"guessNumber/utils"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
	Answer string
	Status playerStatusType.PlayerStatusType
	GameId *string
}

// 監聽 player.Socket.ReadMessage()
func (player *Player) Read(gameServer *GameServer) {
	defer func() {
		_ = player.Socket.Close()
	}()
	for {
		_, message, err := player.Socket.ReadMessage()
		if err != nil {
			gameServer.Unregister <- player
			_ = player.Socket.Close()
			break
		}

		messageHandler(player, gameServer, message)

	}
}

// 監聽 Player.Send
func (player *Player) Write(gameServer *GameServer) {
	defer func() {
		_ = player.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-player.Send:
			if !ok {
				_ = player.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = player.Socket.WriteMessage(websocket.TextMessage, message)

		}
	}
}

func messageHandler(player *Player, gameServer *GameServer, message []byte) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("Server Error: ", r)
		}
	}()

	var data utils.Message

	err := json.Unmarshal(message, &data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data.Type)

	switch data.Type {
	case messageType.INIT:
	case messageType.GET_GAMES:

		gameRespData := gameServer.getGames()

		player.Send <- utils.RespMessage(
			messageType.GET_GAMES, gameRespData,
		)

		return
	case messageType.CREATE_GAME:

		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		gameServer.GameNew <- gameId

		return
	case messageType.GET_PLAYERS:
		break
	case messageType.JOIN_GAME:

		gameId := data.Data.(map[string]interface{})["gameId"].(string)

		if game := gameServer.Game[gameId]; game != nil {
			gameServer.Game[gameId].JoinGame(player)

			gameRespData := gameServer.getGames()

			jsonData := utils.RespMessage(messageType.GET_GAMES, gameRespData)
			gameServer.SendInLobbyPlayers(jsonData)
		}

		return
	case messageType.PLAYING:

		number := data.Data.(map[string]interface{})["value"].(string)

		game := gameServer.Game[*player.GameId]

		ok := utils.ValidateNumber(number)
		// 輸入無效值
		if !ok {
			player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
				Code:    1001,
				Message: "invalid input",
			})
		}

		game.GameHandler(gameServer, number, player)

	case messageType.DELETE_GAME:
		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		fmt.Println(gameId)
		// delete(gameServer.Game, gameId)
	}

}
