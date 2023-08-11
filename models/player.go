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
	Answer *string
	Status playerStatusType.PlayerStatusType
	Game   *Game
}

// 監聽 player.Socket.ReadMessage()
func (player *Player) Read(gameServer *GameServer) {
	defer func() {
		gameServer.Unregister <- player
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
		gameServer.Unregister <- player
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
		if player.Answer != nil {
			gameId := data.Data.(map[string]interface{})["gameId"].(string)

			if game := gameServer.Game[gameId]; game != nil {
				player.JoinGame(game)

				gameRespData := gameServer.getGames()

				jsonData := utils.RespMessage(messageType.GET_GAMES, gameRespData)
				gameServer.SendInLobbyPlayers(jsonData)
			}
		}
	case messageType.INPUT_GAMEANSWER:
		gameAnswer := data.Data.(map[string]interface{})["gameAnswer"].(string)
		player.Answer = &gameAnswer
	case messageType.PLAYING:

		var ok bool

		number := data.Data.(map[string]interface{})["value"].(string)

		if player.Game == nil {
			return
		}

		game := player.Game

		if game != nil {
			return
		}

		ok = utils.ValidateNumber(number)
		// 輸入無效值
		if !ok {
			player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
				Code:    1001,
				Message: "invalid input",
			})
		}

		game.GameHandler(gameServer, number, player)
	case messageType.LEAVE_GAME:
		if player.Game != nil {
			game := player.Game

			player.LeaveGame(game)

			if len(game.Players) <= 0 {
				gameServer.GameEnd <- game
			}
			gameRespData := gameServer.getGames()

			gameServer.SendInLobbyPlayers(utils.RespMessage(
				messageType.GET_GAMES, gameRespData,
			))
		}
	case messageType.DELETE_GAME:
		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		fmt.Println(gameId)
		// delete(gameServer.Game, gameId)
	}

}

func (player *Player) JoinGame(game *Game) {
	player.Game = game

}
func (player *Player) LeaveGame(game *Game) {
	player.Game = nil
	player.Answer = nil
	player.Status = playerStatusType.INLOBBY
	game.Broadcast <- utils.Resp("Opponent Leave. ")
}
