package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"todo-app/pkg/db"
)

// addTaskHandler обрабатывает POST /api/task.
// Десериализует JSON, валидирует поля, добавляет задачу в БД,
// возвращает {"id":"<число>"} или {"error":"<текст>"}.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. Десериализация JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// 2. Заголовок обязателен
	if task.Title == "" {
		writeJSON(w, errorResponse{Error: "не указан заголовок задачи"})
		return
	}

	// 3. Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// 4. Добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, errorResponse{Error: err.Error()})
		return
	}

	// 5. Возвращаем id как строку (фронтенд ожидает string)
	writeJSON(w, map[string]string{"id": itoa(id)})
}

// checkDate проверяет и при необходимости корректирует task.Date.
// Логика:
//   - пустая дата → сегодня
//   - некорректный формат → ошибка
//   - если указан repeat → валидируем правило через NextDate
//   - если дата < сегодня:
//     · без repeat → подставляем сегодня
//     · с repeat  → подставляем ближайшую следующую дату
func checkDate(task *db.Task) error {
	now := time.Now()

	// Пустая дата — берём сегодня
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	// Парсим дату задачи
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("дата указана в неверном формате, ожидается YYYYMMDD")
	}

	// Если задано правило повторения — всегда проверяем его корректность
	// и заодно получаем следующую подходящую дату.
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Если дата задачи не позже сегодняшней — нужна корректировка
	if !t.After(truncateToDay(now)) {
		if task.Repeat == "" {
			// Без правила повторения — ставим сегодня
			task.Date = now.Format(DateFormat)
		} else {
			// С правилом — берём вычисленную следующую дату
			task.Date = next
		}
	}

	return nil
}
