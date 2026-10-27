package main

import (
	"example/speedrun/db"
	"example/speedrun/frontend"
)

func main() {
	db.Database()
	frontend.GenerateTable()
}
