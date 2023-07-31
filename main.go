package main

import (
	"fmt"
	"gin-practice/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	fmt.Println("Starting application...")

	go services.GameServerStart()

	r.GET("/ws/:number", services.GameHandler)

	r.GET("/create/game/:gameId", services.CreateGame)
	r.GET("/join/game/:gameId/:playerId", services.CreateGame)
	r.GET("/list/games", services.GetGames)
	r.GET("/list/players", services.GetOnlinePlayers)
	r.GET("/delete/game/:gameId", services.DeleteGame)

	fmt.Println("chat server start.....")

	r.Run(":8448")
}
