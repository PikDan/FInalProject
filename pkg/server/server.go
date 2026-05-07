package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"todo-app/pkg/api"
	"todo-app/pkg/db"
)

const (
	defaultPort   = 7540
	defaultDBFile = "scheduler.db"
	webDir        = "./web"
)

// Start запускает HTTP-сервер
func Start() error {
	// Инициализируем базу данных
	dbFile := getDBFile()
	if err := db.Init(dbFile); err != nil {
		return err
	}
	log.Printf("database initialized: %s", dbFile)

	//Регистрируем API-обработчики
	api.Init()

	port := getPort()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("starting on http://localhost:%d", port)
	log.Printf("serving static files from: %s", webDir)

	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

// getPort определяет порт для запуска сервера
func getPort() int {
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
			log.Printf("Port taken from TODO_PORT: %d", port)
			return port
		}
		log.Printf("Invalid TODO_PORT: %s. Using default port %d", portStr, defaultPort)
	}
	return defaultPort
}

// getDBFile возвращает путь к файлу базы данных.
// Приоритет: переменная TODO_DBFILE, затем значение по умолчанию.
func getDBFile() string {
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		log.Printf("DB path taken from TODO_DBFILE: %s", path)
		return path
	}
	return defaultDBFile
}

// $env:TODO_PORT=8080; go run .
// //другой порт можно
