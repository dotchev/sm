package main

import (
	"os"

	"github.com/dotchev/sm/rest"
)

func main() {
	router := rest.SMHandler()
	router.Run(listenAddr())
}

func listenAddr() (addr string) {
	addr = os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = "localhost:3000"
	}
	return
}
