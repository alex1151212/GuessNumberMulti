package models

import messageType "guessNumber/enum/message"

type Message struct {
	Data interface{}             `json:"data,omitempty"`
	Type messageType.MessageType `json:"type,omitempty"`
}
