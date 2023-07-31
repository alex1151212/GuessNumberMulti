package main

import (
	"fmt"
	"gin-practice/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	fmt.Println("Starting application...")

	go services.Init()

	r.GET("/ws/:number", services.GameHandler)

	fmt.Println("chat server start.....")

	r.Run(":8448")
}
