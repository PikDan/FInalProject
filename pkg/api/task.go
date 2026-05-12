package api

import (
	"encoding/json"
	"net/http"

	"todo-app/pkg/db" // замените на ваш реальный module path
)

// getTaskHandler обрабатывает GET /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	if task.ID == "" {
		writeJSON(w, errorResponse{Error: "не указан идентификатор"})
		return
	}

	if task.Title == "" {
		writeJSON(w, errorResponse{Error: "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// При успехе возвращаем пустой JSON {}
	writeJSON(w, map[string]any{})
}

// deleteTaskHandler обрабатывает DELETE /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}
