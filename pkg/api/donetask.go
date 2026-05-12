package api

import (
	"net/http"
	"time"

	"todo-app/pkg/db" // замените на ваш реальный module path
)

// doneTaskHandler обрабатывает POST /api/task/done?id=<id>.
// Для одноразовой задачи (repeat пустой) — удаляет её.
// Для периодической — вычисляет следующую дату и обновляет её в БД.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, errorResponse{Error: "не указан идентификатор"})
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// Одноразовая задача — просто удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	// Периодическая задача — вычисляем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// Обновляем только дату в БД
	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}
