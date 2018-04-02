package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestFruits(t *testing.T) {
	handler := SMHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/ping").
		Expect().
		Status(http.StatusOK)
}
