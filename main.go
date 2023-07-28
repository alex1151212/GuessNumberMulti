package main

import (
	"fmt"
	"gin-practice/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()

	fmt.Println("Starting application...")

	go services.Init()

	r.GET("/ws/:number", services.GameHandler)

	fmt.Println("chat server start.....")

	r.Run(":8448")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
