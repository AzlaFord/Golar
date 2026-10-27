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
		description TEXT UNIQUE,
		done INTEGER CHECK (done IN (0,1))
		)`)
	return err
}

func (r *TaskRepository) InsertTask(task Task) error {
	_, err := r.Db.Exec(`INSERT INTO tasks (description,done) VALUES(?,?)`, task.Description, task.Done)

	return err
}

func (r *TaskRepository) ReadTasks() ([]Task, error) {

	rows, err := r.Db.Query(`SELECT * FROM tasks`)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Description, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) MarkDoneDb(id string) error {
	_, err := r.Db.Exec(`UPDATE tasks SET Done=1 WHERE id=?`, id)
	if err != nil {
		return err
	}

	return nil
}
