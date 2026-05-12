package db

import (
	"fmt"
	"time"
)

// Task представляет задачу в таблице scheduler.
// Теги json используются при сериализации/десериализации в API-слое.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает id новой записи.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Tasks возвращает список ближайших задач, отсортированных по дате.
// limit — максимальное количество записей.
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat
	          FROM scheduler
	          ORDER BY date
	          LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Инициализируем пустым слайсом, чтобы JSON вернул [] а не null
	tasks := make([]*Task, 0)

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SearchTasks ищет задачи по подстроке в заголовке или комментарии.
// Если search соответствует формату 02.01.2006 — ищет задачи на конкретную дату.
func SearchTasks(search string, limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	// Проверяем, не является ли строка поиска датой в формате 02.01.2006
	if t, err := time.Parse("02.01.2006", search); err == nil {
		// Поиск по конкретной дате — конвертируем в формат 20060102
		date := t.Format("20060102")
		query := `SELECT id, date, title, comment, repeat
		          FROM scheduler
		          WHERE date = ?
		          ORDER BY date
		          LIMIT ?`

		rows, err := DB.Query(query, date, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				return nil, err
			}
			tasks = append(tasks, &task)
		}
		return tasks, rows.Err()
	}

	// Иначе — поиск по подстроке в title или comment
	like := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat
	          FROM scheduler
	          WHERE title LIKE ? OR comment LIKE ?
	          ORDER BY date
	          LIMIT ?`

	rows, err := DB.Query(query, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

// GetTask возвращает задачу по идентификатору.
func GetTask(id string) (*Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat
	          FROM scheduler WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}
	return &t, nil
}

// UpdateTask обновляет задачу в таблице scheduler по id.
// Возвращает ошибку, если задача с таким id не найдена.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
	          SET date = ?, title = ?, comment = ?, repeat = ?
	          WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// DeleteTask удаляет задачу по идентификатору.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// UpdateDate обновляет только дату задачи по идентификатору.
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
