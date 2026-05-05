package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPort = 7540
	webDir      = "./web"
)

// Start запускает HTTP-сервер
func Start() error {
	port := getPort()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("starting on http://localhost:%d", port)
	log.Printf("serving static files from: %s", webDir)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

// getPort определяет порт для запуска сервера
func getPort() int {
	// Приоритет: переменная TODO_PORT
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
			log.Printf("Port taken from TODO_PORT: %d", port)
			return port
		}
		log.Printf("Invalid TODO_PORT: %s. Using default port %d", portStr, defaultPort)
	}

	return defaultPort
}
