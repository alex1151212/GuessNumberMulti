package models

import (
	"encoding/json"
	gameStatusType "guessNumber/enum/gameStatus"
	messageType "guessNumber/enum/message"
	playerStatusType "guessNumber/enum/playerStatus"
	"guessNumber/utils"
)

type GameServer struct {
	Players map[string]*Player
	Game    map[string]*Game

	GameNew    chan string
	GameEnd    chan *Game
	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan []byte
}

func (gameServer *GameServer) Init() {
	for {
		select {
		case conn := <-gameServer.Register:

			type playerData struct {
				Id string
			}

			jsonData, _ := json.Marshal(&playerData{Id: conn.Id})

			gameServer.SendPlayer(jsonData, conn)

			gameServer.Players[conn.Id] = conn

			conn.Status = playerStatusType.INLOBBY

			gameServer.Broadcast <- utils.Resp("A new socket has connected. ")

		case conn := <-gameServer.Unregister:

			if _, ok := gameServer.Players[conn.Id]; ok {

				gameId := conn.GameId

				if gameId != nil {
					gameServer.Game[*conn.GameId].LeaveGame(conn)
				}

				close(conn.Send)
				delete(gameServer.Players, conn.Id)

				conn.Socket.Close()
				conn = nil

			}

		case game := <-gameServer.GameNew:
			gameServer.createGame(game)

			gameRespData := gameServer.getGames()
			jsonData := utils.RespMessage(messageType.GET_GAMES, gameRespData)

			gameServer.SendInLobbyPlayers(jsonData)

		case game := <-gameServer.GameEnd:
			gameServer.deleteGame(*game.Id)

			gameRespData := gameServer.getGames()
			jsonData := utils.RespMessage(messageType.GET_GAMES, gameRespData)

			gameServer.SendInLobbyPlayers(jsonData)

		case message := <-gameServer.Broadcast:

			for _, conn := range gameServer.Players {
				conn.Send <- message
			}

		}
	}
}

// 發送給指定遊戲內的指定玩家
func (gameServer *GameServer) SendGamePlayer(message []byte, player *Player) {
	for _, conn := range gameServer.Game[*player.GameId].Players {
		if conn == player {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的指定玩家
func (gameServer *GameServer) SendPlayer(message []byte, player *Player) {
	player.Send <- message
}

// 發送給伺服器內閒置在大廳的所有玩家
func (gameServer *GameServer) SendInLobbyPlayers(message []byte) {
	for _, player := range gameServer.Players {
		if player.Status == playerStatusType.INLOBBY {
			gameServer.SendPlayer(message, player)
		}
	}
}

func (gameServer *GameServer) getGames() map[string]*utils.GameRoomRespType {

	var gameList = make(map[string]*utils.GameRoomRespType)
	for k, v := range gameServer.Game {
		gameList[k] = &utils.GameRoomRespType{
			Id:           k,
			PlayerAmount: len(v.Players),
		}
	}
	for _, player := range gameServer.Players {
		if player.Status == playerStatusType.INLOBBY {
			gameServer.SendPlayer(utils.RespMessage(messageType.GET_GAMES, gameList), player)
		}
	}

	return gameList
}

func (gameServer *GameServer) createGame(gameId string) {

	game, ok := gameServer.Game[gameId]

	if !ok {
		gameServer.Game[gameId] = &Game{
			Id:          &gameId,
			Players:     make(map[string]*Player),
			CurrentTurn: nil,
			Winner:      nil,
			Status:      gameStatusType.WAITING,
			Broadcast:   make(chan []byte, 1),
		}
		game = gameServer.Game[gameId]
		game.Init()
	}

	var gameList = make(map[string]*utils.GameRoomRespType)
	for k, v := range gameServer.Game {
		gameList[k] = &utils.GameRoomRespType{
			Id:           k,
			PlayerAmount: len(v.Players),
		}
	}

	for _, player := range gameServer.Players {
		if player.Status == playerStatusType.INLOBBY {
			gameServer.SendPlayer(utils.RespMessage(messageType.GET_GAMES, gameList), player)
		}
	}
}

func (gameServer *GameServer) deleteGame(gameId string) {
	delete(gameServer.Game, gameId)
}
