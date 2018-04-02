package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Platform struct {
	Id          string
	Name        string
	Type        string
	Description string
}

type Platforms struct {
}

func (platforms *Platforms) List(c *gin.Context) {
	list := []Platform{
		{"123-abc", "cf-eu-10", "cloudfoundry", "CF in EU"},
		{"123-xyz", "k8s-05", "kubernetes", "K8S in US"},
	}
	c.JSON(http.StatusOK, gin.H{"platforms": list})
}

func (platforms *Platforms) Create(c *gin.Context) {

}

func (platforms *Platforms) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"platforms": Platform{
		"123-abc", "cf-eu-10", "cloudfoundry", "CF in EU"}})
}

func (platforms *Platforms) Update(c *gin.Context) {

}

func (platforms *Platforms) Delete(c *gin.Context) {

}

func (platforms *Platforms) Register(router gin.IRouter) {
	router.GET("/", platforms.List)
	router.POST("/", platforms.Create)
	router.GET("/:id", platforms.Get)
	router.PATCH("/:id", platforms.Update)
	router.DELETE("/:id", platforms.Delete)
}
