package models

import (
	"encoding/json"
	"fmt"
	gameStatusType "guessNumber/enum/gameStatus"
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

// 傳送訊息給玩家
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
	case messageType.CREATE_GAMES:

		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		gameServer.GameNew <- gameId

		return
	case messageType.GET_PLAYERS:
		break
	case messageType.JOIN_GAME:

		gameId := data.Data.(map[string]interface{})["gameId"].(string)

		if game := gameServer.Game[gameId]; game != nil {
			gameServer.Game[gameId].Join <- player

			gameRespData := gameServer.getGames()
			for _, v := range gameRespData {
				fmt.Println(v)
			}

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

		respA, respB := game.GameHandler(number, player)

		if game.Status == gameStatusType.NORMAL_END {
			respMessage := utils.RespMessage(
				messageType.NORMAL_END, &utils.GameEndRespType{
					GameId:     *game.Id,
					GameStatus: gameStatusType.NORMAL_END,
					Winner:     game.Winner.Id,
				},
			)

			gameServer.Game[*game.Id].Broadcast <- respMessage
			gameServer.GameEnd <- game
			return
		}

		respMessage := utils.RespMessage(
			messageType.PLAYING, &utils.PlayingDataType{
				Resp: utils.PlayingRespType{
					A: respA,
					B: respB,
				},
				Round: game.CurrentTurn.Id,
			},
		)

		if game.CurrentTurn == player {
			for _, gamePlayer := range game.Players {
				if player != gamePlayer {
					game.CurrentTurn = gamePlayer
				}
			}
			gameServer.Game[*game.Id].Broadcast <- respMessage
			return
		}
	case messageType.DELETE_GAME:
		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		fmt.Println(gameId)
		// delete(gameServer.Game, gameId)
	}

}
