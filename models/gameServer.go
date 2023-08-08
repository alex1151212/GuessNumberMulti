package models

import (
	"encoding/json"
	"fmt"
	messageType "guessNumber/enum/message"
	playerStatusType "guessNumber/enum/playerStatus"
	"guessNumber/utils"
)

type GameServer struct {
	Players    map[string]*Player
	Game       map[string]*Game
	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan []byte
}

func (gameServer *GameServer) Init() {
	for {
		select {
		case conn := <-gameServer.Register:

			gameServer.SendPlayers(utils.Resp("A new socket has connected. "), conn)

			type PlayerData struct {
				Id string
			}

			jsonData, _ := json.Marshal(&PlayerData{Id: conn.Id})
			gameServer.SendPlayer(jsonData, conn)

			gameServer.Players[conn.Id] = conn

			conn.Status = playerStatusType.INLOBBY

		case conn := <-gameServer.Unregister:

			if _, ok := gameServer.Players[conn.Id]; ok {
				close(conn.Send)
				delete(gameServer.Players, conn.Id)
				// TODO 通知對手玩家已斷線
				// gameServer.SendGamePlayers(conn.GameId, utils.Resp("socket has disconnected. "), conn)
			}

		// Broadcast
		case message := <-gameServer.Broadcast:

			fmt.Println(message)
			// for conn := range gameServer.Players {
			// 	select {
			// 	case conn.Send <- message:
			// 	default:
			// 		close(conn.Send)
			// 		delete(gameServer.Players, conn)
			// 	}
			// }
		}
	}
}

// 發送給指定遊戲內的所有玩家
func (gameServer *GameServer) SendGamePlayers(gameId string, message []byte, ignore *Player) {
	for _, conn := range gameServer.Game[gameId].Players {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 發送給指定遊戲內的指定玩家
func (gameServer *GameServer) SendGamePlayer(message []byte, player *Player) {
	for _, conn := range gameServer.Game[player.GameId].Players {
		if conn == player {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的所有玩家
func (gameServer *GameServer) SendPlayers(message []byte, ignore *Player) {
	for _, conn := range gameServer.Players {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的指定玩家
func (gameServer *GameServer) SendPlayer(message []byte, player *Player) {
	player.Send <- message
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
