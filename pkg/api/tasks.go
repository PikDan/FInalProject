package api

import (
	"net/http"

	"todo-app/pkg/db" // замените на ваш реальный module path
)

const tasksLimit = 50

// tasksResp — JSON-ответ со списком задач.
type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET /api/tasks и GET /api/tasks?search=...
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	if search != "" {
		tasks, err = db.SearchTasks(search, tasksLimit)
	} else {
		tasks, err = db.Tasks(tasksLimit)
	}

	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, tasksResp{Tasks: tasks})
}
