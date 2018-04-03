package rest

import (
	"net/http"
	"reflect"

	"github.com/dotchev/sm/model"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

type Storage interface {
	AddPlatform(platform *model.Platform) error
	GetPlatforms() (platforms []model.Platform, err error)
}

type Platforms struct {
	storage Storage
}

func NewPlatforms(storage Storage) *Platforms {
	return &Platforms{storage}
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
}

func (platforms *Platforms) List(c *gin.Context) {
	list, err := platforms.storage.GetPlatforms()
	if err != nil {
		sendError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (platforms *Platforms) Create(c *gin.Context) {
	type Platform struct {
		Id          string `json:"id"`
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
	}
	var platform Platform
	if err := c.ShouldBind(&platform); err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	if platform.Id == "" {
		platform.Id = uuid.NewV4().String()
	}

	err := platforms.storage.AddPlatform((*model.Platform)(&platform))
	if err != nil {
		sendError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, &platform)
}

func (platforms *Platforms) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"platforms": model.Platform{
		"123-abc", "cf-eu-10", "cloudfoundry", "CF in EU"}})
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
