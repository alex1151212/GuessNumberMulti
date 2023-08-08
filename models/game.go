package models

import (
	"fmt"
	gameStatusType "guessNumber/enum/gameStatus"
	"strings"

	"github.com/gorilla/websocket"
)

type Game struct {
	Id          string
	Players     []*Player
	CurrentTurn *Player
	Winner      *Player
	Status      gameStatusType.GameStatusType
}

var GAEM_MAX_PLAYER_AMOUNT = 2

func (game *Game) Init() {
	game.Status = gameStatusType.START
	game.CurrentTurn = game.Players[0]
}

func PlayerAmountHandler(game Game, player *Player) {
	if len(game.Players)+1 > 2 {
		_ = player.Socket.WriteMessage(websocket.TextMessage, []byte("Room is full"))
		player.Socket.Close()
		return
	}
}

// 遊戲邏輯
func (game *Game) gameResponse(number string, c *Player) (respA string, respB string) {
	var answer map[int]string = make(map[int]string)
	var a, b int
	if game.Status == gameStatusType.START {
		for _, client := range game.Players {
			if client != c {
				strAnswerList := strings.Split(client.Answer, "")
				for index, item := range strAnswerList {
					answer[index] = item
				}
			}
		}
		strNumberList := strings.Split(number, "")
		for index, item := range strNumberList {
			if answer[index] == item {

				a += 1
			} else {
				for _, value := range answer {
					if value == item {
						b += 1
					}
				}
			}
		}

		if a == 4 {
			game.Winner = c
			game.Status = gameStatusType.NORMAL_END
		}
		respA = fmt.Sprintf("%d", a)
		respB = fmt.Sprintf("%d", b)

	}

	return
}
