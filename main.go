package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := SMHandler()

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = "localhost:3000"
	}
	r.Run(addr)
}

func SMHandler() *gin.Engine {
	router := gin.Default()

	platforms := Platforms{}
	platforms.Register(router.Group("/v1/platforms"))

	return router
}
