package main

import (
	"log"
	"todo-app/pkg/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
