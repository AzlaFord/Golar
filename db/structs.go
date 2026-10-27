package db

import "database/sql"

type TaskRepository struct {
	Db *sql.DB
}

type Task struct {
	Id          int
	Description string
	Done        bool
}

func (r *TaskRepository) CreateTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS tasks ( id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT,
		done BOOL
		)`)
	return err
}
