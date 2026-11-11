package main

import (
	tea "charm.land/bubbletea/v2"
	"database/sql"
	// "example/speedrun/db"
	"example/speedrun/frontend"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

func main() {

	DbConnection, err := sql.Open("sqlite", "./base.db")

	if err != nil {
		fmt.Println(err)
	}

	defer DbConnection.Close()
	// db.Database(DbConnection)
	tasks, err := frontend.LoadTasks(DbConnection)
	if err != nil {
		fmt.Println(err)
	}

	m := frontend.GenerateTable(DbConnection, tasks)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
