package main

import (
	"fmt"
	"guessNumber/services"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	// corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.MaxAge = 12 * time.Hour
	r.Use(cors.New(corsConfig))

	fmt.Println("Starting application...")

	go services.GameServerStart()

	r.GET("/ws/:number", services.GameHandler)

	r.GET("/create/game/:gameId", services.CreateGame)
	r.GET("/join/game/:gameId/:playerId", services.JoinGame)
	r.GET("/list/games", services.GetGames)
	r.GET("/list/players", services.GetOnlinePlayers)
	r.GET("/delete/game/:gameId", services.DeleteGame)

	fmt.Println("Game server start.....")

	go enablePprofServer()

	r.Run(":8448")
}

func enablePprofServer() {
	http.ListenAndServe("localhost:6060", nil)
}
