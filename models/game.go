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
	Join      chan *Player
	Leave     chan *Player
	Broadcast chan []byte
}

var GAEM_MAX_PLAYER_AMOUNT = 2

func (game *Game) Init(gameServer *GameServer) {
	game.Status = gameStatusType.WAITING

	go game.gamePlayerHandler(gameServer)
}

func (game *Game) gamePlayerHandler(gameServer *GameServer) {
	for {
		select {
		case message := <-game.Broadcast:
			for _, player := range game.Players {
				player.Send <- message
			}
		case player := <-game.Join:

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

			switch len(game.Players) {
			case 1:
				player.Status = playerStatusType.WAITING_START
			case 2:
				game.startGame()
				for _, player := range game.Players {
					player.Send <- utils.RespMessage(messageType.GAME_START, nil)
					player.Status = playerStatusType.PLAYING
				}
			default:
				player.Send <- utils.RespErrorMessage(utils.ErrorRespType{
					Code:    1003,
					Message: "The Game Room is Full",
				})
				return
			}

			gameRespData := gameServer.getGames()

			jsonData := utils.RespMessage(messageType.GET_GAMES, gameRespData)
			gameServer.SendInLobbyPlayers(jsonData)

		case player := <-game.Leave:

			delete(game.Players, player.Id)

			if len(game.Players) <= 0 {
				gameServer.GameEnd <- game
			}

			gameRespData := gameServer.getGames()

			gameServer.SendInLobbyPlayers(utils.RespMessage(
				messageType.GET_GAMES, gameRespData,
			))
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

			for _, player := range game.Players {
				player.LeaveGame(game)
				game.Leave <- player
			}

			return
		}

	}
}

func calculateGameResult(game *Game, player *Player, number string) (a int, b int, valid bool) {
	var answer map[int]string = make(map[int]string)

	// 取得對手答案
	for _, client := range game.Players {
		if client != player {
			strAnswerList := strings.Split(*client.Answer, "")
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
