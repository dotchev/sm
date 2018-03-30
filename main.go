package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := NewHandler()

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = "localhost:3000"
	}
	r.Run(addr)
}

func NewHandler() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}
