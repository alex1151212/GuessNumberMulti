package utils

import (
	"encoding/json"
	gameStatusType "guessNumber/enum/gameStatus"
	messageType "guessNumber/enum/message"
)

type Message struct {
	Data interface{}             `json:"data,omitempty"`
	Type messageType.MessageType `json:"type,omitempty"`
}
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

func Resp(message string) []byte {
	jsonData, _ := json.Marshal(&Message{Data: message})
	return jsonData
}

func RespMessage(messageType messageType.MessageType, data interface{}) []byte {
	var jsonData []byte
	if data == nil {
		jsonData, _ = json.Marshal((&Message{Type: messageType}))
	} else {
		jsonData, _ = json.Marshal((&Message{Data: data, Type: messageType}))
	}
	return jsonData
}

func RespErrorMessage(data ErrorRespType) []byte {
	jsonData := RespMessage(messageType.ERROR, data)
	return jsonData
}
