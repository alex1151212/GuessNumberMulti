package models

import (
	"encoding/json"
	"fmt"
	messageType "gin-practice/enum/message"
	playerStatusType "gin-practice/enum/playerStatus"
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
				// gameServer.SendGamePlayers(conn.GameId, utils.Resp("socket has disconnected. "), conn)
			}

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

func (gameServer *GameServer) getGames() []byte {
	type gameResponse struct {
		Id           string `json:"id"`
		PlayerAmount int    `json:"playerAmount"`
	}

	var gameList = make([]*gameResponse, 0)
	for _, v := range gameServer.Game {
		gameList = append(gameList, &gameResponse{
			Id:           v.Id,
			PlayerAmount: len(v.Players),
		})
	}
	gameRespData, err := json.Marshal(gameList)
	if err != nil {
		fmt.Println(gameRespData)
	}
	return gameRespData
}

func (gameServer *GameServer) joinGame(gameId string, playerId string) {
	game := gameServer.Game[gameId]

	for _, player := range game.Players {
		if player.Id == playerId {
			return
		}
	}
	game.Players = append(game.Players, gameServer.Players[playerId])
	gameServer.Players[playerId].GameId = gameId

	if len(game.Players) == 2 {
		game.Init()
		gameServer.SendGamePlayers(gameId, utils.RespMessage(&utils.Message{
			Type: messageType.GAME_START,
		}), nil)

		for _, player := range game.Players {
			player.Status = playerStatusType.PLAYING
		}
	} else {
		gameServer.Players[playerId].Status = playerStatusType.WAITINGSTART
	}

}
