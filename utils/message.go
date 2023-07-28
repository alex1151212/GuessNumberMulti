package utils

import "encoding/json"

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	// Test      string `json:"test,omitempty"`
}

// jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected. "})
// game.Send(jsonMessage, conn)

func FormatToJson(message string) []byte {
	jsonData, _ := json.Marshal(&Message{Content: message})
	return jsonData
}
