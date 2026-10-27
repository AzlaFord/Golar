package db

import (
	"database/sql"
	"fmt"
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
	err = taskRepository.InsertTask(Task{1, "Vreau sa fac ceva", false})

	if err != nil {
		log.Fatal("Error inserting inorders table", err)
	}

	tasks, err := taskRepository.ReadTasks()
	if err != nil {
		log.Fatal("Error inserting inorders table", err)
	}
	fmt.Println(tasks)
}
