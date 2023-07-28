package models

import (
	"github.com/gorilla/websocket"
)

type Game struct {
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
