package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// schema  SQL-команд для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)      NOT NULL DEFAULT "",
    title   VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT         NOT NULL DEFAULT "",
    repeat  VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

// открывает базу данных и создаёт таблицу с индексом.
func Init(dbFile string) error {
	// Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем файл базы данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверяем соединение
	if err = DB.Ping(); err != nil {
		return err
	}

	// Если файл не существовал — создаём таблицу и индекс
	if install {
		if _, err = DB.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}
