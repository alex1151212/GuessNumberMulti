package models

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Player struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
	Answer string
	// TODO 實作房間機制
	RoomId string
	// 回合制
	// 1. isMyTurn判斷
	// 2. mutex加鎖
}

func (c *Player) Read() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			// manager.unregister <- c
			_ = c.Socket.Close()
			break
		}
		fmt.Println(message)
		//		messageToStr := string(message)
		//		strSplit := strings.Split(messageToStr, "")
		//		isValid := true
		//		for _, str := range strSplit {
		//			_, err := strconv.Atoi(str)
		//			if err != nil {
		//				isValid = false
		//			}
		//		}
		//		if c.isMyTurn && isValid {
		//			for client := range manager.clients {
		//				if client == c {
		//					client.isMyTurn = false
		//					res := gameResponse(string(message), c)
		//					jsonMessage, _ := json.Marshal(&Message{Content: "/A Your guess number is: " + string(message)})
		//					manager.SendMySelf(jsonMessage, c)
		//					jsonMessage, _ = json.Marshal(&Message{Content: "/A Response: " + res})
		//					manager.SendMySelf(jsonMessage, c)
		//				} else {
		//					client.isMyTurn = true
		//					jsonMessage, _ := json.Marshal(&Message{Content: "/A User guess your number is: " + string(message)})
		//					manager.send(jsonMessage, c)
		//					jsonMessage, _ = json.Marshal(&Message{Content: "/A It's your turn. "})
		//					manager.send(jsonMessage, c)
		//				}
		//			}
		//		} else if c.isMyTurn {
		//			jsonMessage, _ := json.Marshal(&Message{Content: "/A Your guess number is not valid: " + string(message)})
		//			manager.SendMySelf(jsonMessage, c)
		//		} else {
		//			jsonMessage, _ := json.Marshal(&Message{Content: "/A Its not your turn "})
		//			manager.SendMySelf(jsonMessage, c)
		//		}
	}
}

func (c *Player) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.Socket.WriteMessage(websocket.TextMessage, message)

		}
	}
}
