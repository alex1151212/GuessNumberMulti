package models

import gameStatusType "guessNumber/enum/gameStatus"

type JoinGameDataType struct {
	PlayerId string `json:"playerId,omitempty"`
	GameId   string `json:"gameId,omitempty"`
}
type CreateGameDataType struct {
	GameId string `json:"gameId,omitempty"`
}
type PlayingDataType struct {
	Value string          `json:"value,omitempty"`
	Resp  PlayingRespType `json:"resp,omitempty"`
	Round string          `json:"round,omitempty"`
}
type PlayingRespType struct {
	A string `json:"a"`
	B string `json:"b"`
}
type GameRoomRespType struct {
	Id           string `json:"id"`
	PlayerAmount int    `json:"playerAmount"`
}
type GameEndRespType struct {
	GameId     string                        `json:"gameId"`
	Winner     string                        `json:"winner"`
	GameStatus gameStatusType.GameStatusType `json:"gameStatus"`
}
type ErrorRespType struct {
	Code    int
	Message string
}
