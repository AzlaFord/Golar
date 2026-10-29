package main

import (
	"database/sql"
	"example/speedrun/db"
	"example/speedrun/frontend"
	"fmt"
	_ "modernc.org/sqlite"
)

func main() {

	DbConnection, err := sql.Open("sqlite", "./base.db")

	if err != nil {
		fmt.Println(err)
	}

	defer DbConnection.Close()
	db.Database(DbConnection)
	frontend.GenerateTable(DbConnection)
}
