package db

import (
	"database/sql"
	"log"
)

func Database(dbConnection *sql.DB) {

	taskRepository := &TaskRepository{Db: dbConnection}

	err := taskRepository.CreateTable()
	if err != nil {
		log.Fatal("Error creating orders table", err)
	}
}

func AddTask(dbConnection *sql.DB, task string) error {

	taskRepository := &TaskRepository{Db: dbConnection}
	err := taskRepository.InsertTask(task)
	if err != nil {
		return err
	}
	return nil
}

func GetList(dbConnection *sql.DB) ([]Task, error) {

	taskRepository := &TaskRepository{Db: dbConnection}
	tasks, err := taskRepository.ReadTasks()
	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func ClearDoneTasks(dbConnection *sql.DB) error {
	TaskRepository := &TaskRepository{Db: dbConnection}
	err := TaskRepository.DeleteDoneTasks()
	if err != nil {
		return err
	}
	return nil
}
func MarkDone(id string, dbConnection *sql.DB) error {

	taskRepository := &TaskRepository{Db: dbConnection}
	err := taskRepository.MarkDoneDb(id)
	if err != nil {
		return err
	}
	return nil
}
