package rest

import (
	"github.com/dotchev/sm/storage/postgres"
	"github.com/gin-gonic/gin"
)

func SMHandler() *gin.Engine {
	storage, err := postgres.NewStorage(
		"postgres://postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	platforms := NewPlatforms(storage)
	platforms.Register(router.Group("/v1/platforms"))

	return router
}
