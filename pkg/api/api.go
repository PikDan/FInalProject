package api

import "net/http"

const DateFormat = "20060102"

// Init регистрирует все API-обработчики.
func Init() {
	// Публичный эндпоинт — аутентификация
	http.HandleFunc("/api/signin", signinHandler)

	// Вспомогательный эндпоинт — не требует аутентификации
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// Защищённые эндпоинты — оборачиваем в auth()
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}

// taskHandler — единая точка входа для /api/task.
// Маршрутизирует запрос по HTTP-методу.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, errorResponse{Error: "метод не поддерживается"})
	}
}
