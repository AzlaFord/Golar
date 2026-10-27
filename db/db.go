package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func Database() {
	dbConnection, err := sql.Open("sqlite", "./base.db")

	if err != nil {
		log.Fatal(err)
	}

	defer dbConnection.Close()

	taskRepository := &TaskRepository{Db: dbConnection}

	err = taskRepository.CreateTable()
	if err != nil {
		log.Fatal("Error creating orders table", err)
	}
	// err = taskRepository.InsertTask(Task{1, "Altceva sa fac ceva, sa fac", false})

	// if err != nil {
	// 	log.Fatal("Error inserting inorders table", err)
	// }

}
func GetList() ([]Task, error) {

	dbConnection, err := sql.Open("sqlite", "./base.db")

	if err != nil {
		log.Fatal(err)
	}

	defer dbConnection.Close()
	taskRepository := &TaskRepository{Db: dbConnection}
	tasks, err := taskRepository.ReadTasks()
	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func MarkDone(id string) error {

	dbConnection, err := sql.Open("sqlite", "./base.db")

	if err != nil {
		log.Fatal(err)
	}

	defer dbConnection.Close()
	taskRepository := &TaskRepository{Db: dbConnection}
	err = taskRepository.MarkDoneDb(id)
	if err != nil {
		return err
	}
	return nil
}
