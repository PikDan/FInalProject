package api

import "net/http"

const DateFormat = "20060102"

// Init регистрирует API-обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}
