package models

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type Game struct {
	Id          uint
	Players     []*Player
	CurrentTurn *Player
}

var GAEM_MAX_PLAYER_AMOUNT = 2

func (game *Game) Init() {
	game.CurrentTurn = game.Players[0]
}

func PlayerAmountHandler(game Game, player *Player) {
	if len(game.Players)+1 > 2 {
		_ = player.Socket.WriteMessage(websocket.TextMessage, []byte("Room is full"))
		player.Socket.Close()
		return
	}
}

func (game *Game) gameResponse(number string, c *Player) string {
	var answer map[string]string = make(map[string]string)
	var a, b int
	var returnStr string
	for _, client := range game.Players {
		if client != c {
			strAnswerList := strings.Split(client.Answer, "")
			for index, item := range strAnswerList {
				answer[string(index)] = item
			}
		}
	}
	strNumberList := strings.Split(number, "")

	for index, item := range strNumberList {
		if answer[string(index)] == item {
			a += 1
		} else {
			for _, value := range answer {
				if value == item {
					b += 1
				}
			}
		}
	}
	returnStr = fmt.Sprintf("%dA %dB", a, b)

	return returnStr
}
