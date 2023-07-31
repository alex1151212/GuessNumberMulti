package models

type GameServer struct {
	Players    map[*Player]bool
	Game       *Game
	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan []byte
}

// 發送給指定遊戲內的所有玩家
func (gameServer *GameServer) SendGamePlayers(message []byte, ignore *Player) {
	for _, conn := range gameServer.Game.Players {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 發送給指定遊戲內的指定玩家
func (gameServer *GameServer) SendGamePlayer(message []byte, player *Player) {
	for _, conn := range gameServer.Game.Players {
		if conn == player {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的所有玩家
func (gameServer *GameServer) SendPlayers(message []byte, ignore *Player) {
	for conn := range gameServer.Players {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// 發送給伺服器內的指定玩家
func (gameServer *GameServer) SendPlayer(message []byte, player *Player) {
	for conn := range gameServer.Players {
		if conn == player {
			conn.Send <- message
		}
	}
}
