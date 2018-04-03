package main

import (
	"os"

	"github.com/dotchev/sm/postgres"
	"github.com/dotchev/sm/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	router := SMHandler()
	router.Run(listenAddr())
}

func listenAddr() (addr string) {
	addr = os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = "localhost:3000"
	}
	return
}

func SMHandler() *gin.Engine {
	_, err := postgres.NewStorage(
		"postgres://postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	platforms := rest.Platforms{}
	platforms.Register(router.Group("/v1/platforms"))

	return router
}
