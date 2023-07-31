package models

import (
	"fmt"
	"gin-practice/utils"
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

			// gameServer.Players[conn] = true

			gameServer.SendPlayers(utils.FormatToJson("A new socket has connected. "), conn)

		case conn := <-gameServer.Unregister:

			fmt.Println(conn)
			// if _, ok := gameServer.Players[conn]; ok {
			// 	close(conn.Send)
			// 	delete(gameServer.Players, conn)
			// 	gameServer.SendGamePlayers(conn.GameId, utils.FormatToJson("socket has disconnected. "), conn)
			// }

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
	for _, conn := range gameServer.Players {
		if conn == player {
			conn.Send <- message
		}
	}
}
