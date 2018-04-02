package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sm Suite")
}

var _ = Describe("REST", func() {
	var server *httptest.Server
	var e *httpexpect.Expect

	BeforeSuite(func() {
		handler := SMHandler()
		server = httptest.NewServer(handler)
		e = httpexpect.New(GinkgoT(), server.URL)
	})

	AfterSuite(func() {
		server.Close()
	})

	It("/ping should respond with pong", func() {
		e.GET("/ping").
			Expect().
			Status(http.StatusOK).
			JSON().Equal(gin.H{
			"message": "pong",
		})
	})

})
