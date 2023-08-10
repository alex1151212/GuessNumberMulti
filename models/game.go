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

func (game *Game) startGame() {

	game.Status = gameStatusType.START
	game.initCurrentRound()
}

func (game *Game) JoinGame(player *Player) {

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
}
func (game *Game) LeaveGame(player *Player) {
	delete(game.Players, player.Id)
	game.Broadcast <- utils.Resp("Opponent Leave. ")
}

// 遊戲邏輯
func (game *Game) GameHandler(gameServer *GameServer, number string, player *Player) {

	var a, b int
	var valid bool
	var respA, respB string
	if game.Status == gameStatusType.START {

		if game.CurrentTurn == player {

			a, b, valid = calculateGameResult(game, player, number)

			if !valid {
				fmt.Println("inptu invalid")
				return
			}

			respA = fmt.Sprintf("%d", a)
			respB = fmt.Sprintf("%d", b)

			respMessage := utils.RespMessage(
				messageType.PLAYING, &utils.PlayingDataType{
					Resp: utils.PlayingRespType{
						A: respA,
						B: respB,
					},
					Round: game.CurrentTurn.Id,
				},
			)

			game.Broadcast <- respMessage

			for _, gamePlayer := range game.Players {
				if player != gamePlayer {
					game.CurrentTurn = gamePlayer
				}
			}

		} else {
			return
		}

		if a == 4 {
			game.Winner = player
			game.Status = gameStatusType.NORMAL_END

			for _, player := range game.Players {
				player.GameId = nil
			}
		}
		if game.Status == gameStatusType.NORMAL_END {
			respMessage := utils.RespMessage(
				messageType.NORMAL_END, &utils.GameEndRespType{
					GameId:     *game.Id,
					GameStatus: gameStatusType.NORMAL_END,
					Winner:     game.Winner.Id,
				},
			)

			game.Broadcast <- respMessage

			gameServer.GameEnd <- game

			return
		}

	}
}

func calculateGameResult(game *Game, player *Player, number string) (a int, b int, valid bool) {
	var answer map[int]string = make(map[int]string)

	// 取得對手答案
	for _, client := range game.Players {
		if client != player {
			strAnswerList := strings.Split(client.Answer, "")
			for index, item := range strAnswerList {
				answer[index] = item
			}
		}
	}

	charMap := make(map[rune]bool)

	for _, char := range number {
		if charMap[char] {
			valid = false
			return
		}
		charMap[char] = true
	}
	valid = true

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
	return
}
