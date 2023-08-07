package models

import (
	"encoding/json"
	"fmt"
	messageType "gin-practice/enum/message"
	playerStatusType "gin-practice/enum/playerStatus"
	"gin-practice/utils"
	"strconv"
	"strings"

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

func messageHandler(myself *Player, gameServer *GameServer, message []byte) {

	var data utils.Message

	err := json.Unmarshal(message, &data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data.Type)

	switch data.Type {
	case messageType.INIT:
		gameRespData := gameServer.getGames()

		myself.Send <- utils.RespMessage(&utils.Message{
			Type: messageType.GET_GAMES,
			Data: string(gameRespData),
		})
		return
	case messageType.GET_GAMES:

		gameRespData := gameServer.getGames()

		myself.Send <- utils.RespMessage(&utils.Message{
			Type: messageType.GET_GAMES,
			Data: string(gameRespData),
		})

		return
	case messageType.CREATE_GAMES:
		break
	case messageType.GET_PLAYERS:
		break
	case messageType.JOIN_GAME:

		gameId := data.Data.(map[string]interface{})["gameId"].(string)
		playerId := data.Data.(map[string]interface{})["playerId"].(string)

		gameServer.joinGame(gameId, playerId)

		gameRespData := gameServer.getGames()

		for _, player := range gameServer.Players {
			if player.Status == playerStatusType.INLOBBY {
				gameServer.SendPlayer(utils.RespMessage(&utils.Message{
					Type: messageType.GET_GAMES,
					Data: string(gameRespData),
				}), player)
			}
		}

		// myself.Send <- utils.RespMessage(&utils.Message{
		// 	Type: messageType.GET_GAMES,
		// 	Data: string(gameRespData),
		// })

		return
	case messageType.PLAYING:

		number := data.Data.(map[string]interface{})["value"].(string)

		game := gameServer.Game[myself.GameId]

		messageToStr := string(number)
		strSplit := strings.Split(messageToStr, "")
		isValid := true
		for _, str := range strSplit {
			_, err := strconv.Atoi(str)
			if err != nil {
				isValid = false
			}
		}
		respA, respB := game.gameResponse(messageToStr, myself)
		respMessage := utils.RespMessage(&utils.Message{
			Type: messageType.PLAYING,
			Data: &utils.PlayingDataType{
				Resp: utils.PlayingRespType{
					A: respA,
					B: respB,
				},
				Round: game.CurrentTurn.Id,
			},
		})
		if game.CurrentTurn == myself && isValid {
			for _, player := range game.Players {
				if player == myself {
					gameServer.SendGamePlayer(respMessage, player)
				} else {
					gameServer.SendGamePlayer(respMessage, player)
					game.CurrentTurn = player
				}

			}
			return
		} else if !isValid {
			return
			// Your guess number is not valid:
		} else {
			return
			// Its not your turn.
		}

	}

}
