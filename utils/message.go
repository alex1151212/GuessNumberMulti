package utils

import (
	"encoding/json"
	messageType "gin-practice/enum/message"
)

type Message struct {
	Sender    string                  `json:"sender,omitempty"`
	Recipient string                  `json:"recipient,omitempty"`
	Data      interface{}             `json:"data,omitempty"`
	Type      messageType.MessageType `json:"type,omitempty"`
	// Test      string `json:"test,omitempty"`
}
type JoinGameDataType struct {
	PlayerId string `json:"playerId,omitempty"`
	GameId   string `json:"gameId,omitempty"`
}
type PlayingDataType struct {
	Value string          `json:"value,omitempty"`
	Resp  PlayingRespType `json:"resp,omitempty"`
	Round string          `json:"round,omitempty"`
}
type PlayingRespType struct {
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
}

func Resp(message string) []byte {
	jsonData, _ := json.Marshal(&Message{Data: message})
	return jsonData
}

func RespMessage(message *Message) []byte {
	jsonData, _ := json.Marshal(message)
	return jsonData
}
