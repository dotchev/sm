package main

import (
	"os"
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
