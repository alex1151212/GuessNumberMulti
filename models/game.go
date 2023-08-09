package models

import (
	"fmt"
	gameStatusType "guessNumber/enum/gameStatus"
	messageType "guessNumber/enum/message"
	playerStatusType "guessNumber/enum/playerStatus"
	"guessNumber/utils"
	"strings"
)

type Game struct {
	Id          *string
	Players     map[string]*Player
	CurrentTurn *Player
	Winner      *Player
	Status      gameStatusType.GameStatusType

	// 玩家管理頻道
	Leave     chan *Player
	Join      chan *Player
	Broadcast chan []byte
}

var GAEM_MAX_PLAYER_AMOUNT = 2

func (game *Game) Init() {
	game.Status = gameStatusType.WAITING

	go game.gamePlayerHandler()
}

func (game *Game) gamePlayerHandler() {
	for {
		select {
		case player := <-game.Join:
			joinGame(game, player)

			// 開始遊戲判斷
			if len(game.Players) == 2 {
				player.Send <- utils.RespMessage(messageType.GAME_START, nil)
				game.startGame()
				for _, player := range game.Players {
					player.Status = playerStatusType.PLAYING
				}
			} else {
				player.Status = playerStatusType.WAITING_START
			}

		case player := <-game.Leave:

			delete(game.Players, player.Id)
			game.Broadcast <- utils.Resp("Opponent Leave. ")

		case message := <-game.Broadcast:
			for _, player := range game.Players {
				player.Send <- message
			}
		}
	}
}

func (game *Game) initCurrentRound() {
	keys := make([]string, 0, len(game.Players))
	for key := range game.Players {
		keys = append(keys, key)
	}
	game.CurrentTurn = game.Players[keys[0]]
}

func joinGame(game *Game, player *Player) {

	// 遊戲已滿
	if len(game.Players) >= 2 {

		player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
			Code:    1003,
			Message: "The Game Room is Full",
		})
		return
	}

	// 嘗試重複新增同樣的玩家
	for _, gamePlayer := range game.Players {
		if gamePlayer.Id == player.Id {
			player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
				Code:    1002,
				Message: "The Game Already Exist this Player",
			})
			return
		}
	}

	game.Players[player.Id] = player
	player.GameId = game.Id

}

func (game *Game) startGame() {

	game.Status = gameStatusType.START
	game.initCurrentRound()

}

// 遊戲邏輯
func (game *Game) GameHandler(number string, c *Player) (respA string, respB string) {
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

			for _, player := range game.Players {
				player.GameId = nil

			}

		}
		respA = fmt.Sprintf("%d", a)
		respB = fmt.Sprintf("%d", b)

	}

	return
}
