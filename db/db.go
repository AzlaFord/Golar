package db

import (
	"database/sql"
	"log"
	_"github.com/mattn/go-sqlite3"
)

func Database() {
	dbConnection, err := sql.Open("sqlite3", "./base.db")

	if err != nil {
		log.Fatal(err)
	}

	defer dbConnection.Close()

	taskRepository := &TaskRepository{Db: dbConnection}

	err = taskRepository.CreateTable()
	if err != nil {
		log.Fatal("Error creating orders table", err)
	}

}
