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
	GameId string
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
		// TODO WebSocket CreateGame
		// gameId := data.Data.(map[string]interface{})["gameId"].(string)

		break
	case messageType.GET_PLAYERS:
		break
	case messageType.JOIN_GAME:

		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		playerId := data.Data.(map[string]interface{})["playerId"].(string)

		joinGame(gameServer, player, gameId, playerId)

		gameRespData := gameServer.getGames()

		// 對所有閒置在大廳的玩家發送房間更新資料
		for _, player := range gameServer.Players {
			if player.Status == playerStatusType.INLOBBY {
				gameServer.SendPlayer(utils.RespMessage(messageType.GET_GAMES, gameRespData), player)
			}
		}
		return
	case messageType.PLAYING:

		number := data.Data.(map[string]interface{})["value"].(string)

		game := gameServer.Game[player.GameId]

		ok := utils.ValidateNumber(number)
		// 輸入無效值
		if !ok {
			player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
				Code:    1001,
				Message: "invalid input",
			})
		}

		respA, respB := game.gameResponse(number, player)

		if game.Status == gameStatusType.NORMAL_END {
			respMessage := utils.RespMessage(
				messageType.NORMAL_END, &utils.GameEndRespType{
					GameId:     game.Id,
					GameStatus: gameStatusType.NORMAL_END,
					Winner:     game.Winner.Id,
				},
			)
			gameServer.SendGamePlayers(game.Id, respMessage, nil)
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
			gameServer.SendGamePlayers(game.Id, respMessage, nil)
			return
		}

	}

}

func joinGame(gameServer *GameServer, player *Player, gameId string, playerId string) {

	game, ok := gameServer.Game[gameId]

	// 找不到 gameId 的遊戲
	if !ok {
		gameServer.SendPlayer(utils.RespMessage(messageType.GAME_START, nil), player)
		return
	}

	// 遊戲已滿
	if len(game.Players) >= 2 {
		player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
			Code:    1003,
			Message: "The Game Room is Full",
		})
		return
	}

	// 嘗試重複新增同樣的玩家
	for _, player := range game.Players {
		if player.Id == playerId {
			player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
				Code:    1002,
				Message: "The Game Already Exist this Player",
			})
			return
		}
	}

	game.Players = append(game.Players, gameServer.Players[playerId])
	gameServer.Players[playerId].GameId = gameId

	if len(game.Players) == 2 {
		game.Init()
		gameServer.SendGamePlayers(gameId, utils.RespMessage(messageType.GAME_START, nil), nil)

		for _, player := range game.Players {
			player.Status = playerStatusType.PLAYING
		}
	} else {
		gameServer.Players[playerId].Status = playerStatusType.WAITING_START
	}
}
