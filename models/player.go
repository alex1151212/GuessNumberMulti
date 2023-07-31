package models

import (
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
	// TODO 實作房間機制
	GameId string
	// 回合制
	// 1. isMyTurn判斷
	// 2. mutex加鎖
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

		read(player, gameServer, message)

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

func read(myself *Player, gameServer *GameServer, message []byte) {
	game := gameServer.Game[myself.GameId]
	messageToStr := string(message)
	strSplit := strings.Split(messageToStr, "")
	isValid := true
	for _, str := range strSplit {
		_, err := strconv.Atoi(str)
		if err != nil {
			isValid = false
		}
	}
	if game.CurrentTurn == myself && isValid {
		for _, player := range game.Players {
			if player == myself {
				res := game.gameResponse(string(message), myself)
				gameServer.SendGamePlayer(utils.FormatToJson("Your guess number is: "+string(message)), myself)
				gameServer.SendGamePlayer(utils.FormatToJson("Response: "+res), myself)
			} else {
				gameServer.SendGamePlayer(utils.FormatToJson("User guess your number is: "+string(message)), player)
				gameServer.SendGamePlayer(utils.FormatToJson("It's your turn. "+string(message)), player)
				game.CurrentTurn = player
			}

		}
	} else if !isValid {
		gameServer.SendGamePlayer(utils.FormatToJson("Your guess number is not valid: "+string(message)), myself)
	} else {
		gameServer.SendGamePlayer(utils.FormatToJson("Its not your turn. "), myself)
	}
}
