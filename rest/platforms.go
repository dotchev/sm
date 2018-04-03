package rest

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/dotchev/sm/storage"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

type Platforms struct {
	store storage.Storage
}

func NewPlatforms(store storage.Storage) *Platforms {
	return &Platforms{store}
}

func sendError(c *gin.Context, status int, err error) {
	type ErrorReply struct {
		Error       string `json:"error"`
		Description string `json:"description"`
	}
	c.AbortWithStatusJSON(status, &ErrorReply{
		Error:       reflect.TypeOf(err).String(),
		Description: err.Error(),
	})
	fmt.Printf("Error: %v\n", err)
}

func (platforms *Platforms) List(c *gin.Context) {
	type Platform struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}

	list, err := platforms.store.GetPlatforms()
	if err != nil {
		sendError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"platforms": list})
}

func (platforms *Platforms) Create(c *gin.Context) {
	type Platform struct {
		ID          string `json:"id"`
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
	}
	var platform Platform
	if err := c.ShouldBind(&platform); err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	if platform.ID == "" {
		platform.ID = uuid.NewV4().String()
	}

	err := platforms.store.AddPlatform((*storage.Platform)(&platform))
	if err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, &platform)
}

func (platforms *Platforms) Get(c *gin.Context) {
	id := c.Param("id")
	platform, err := platforms.store.GetPlatformByID(id)
	switch {
	case err == storage.ErrNotFound:
		sendError(c, http.StatusNotFound, err)
	case err != nil:
		sendError(c, http.StatusInternalServerError, err)
	default:
		c.JSON(http.StatusOK, platform)
	}
}

func (platforms *Platforms) Update(c *gin.Context) {

}

func (platforms *Platforms) Delete(c *gin.Context) {

}

func (platforms *Platforms) Register(router gin.IRouter) {
	router.GET("", platforms.List)
	router.POST("", platforms.Create)
	router.GET(":id", platforms.Get)
	router.PATCH(":id", platforms.Update)
	router.DELETE(":id", platforms.Delete)
}
