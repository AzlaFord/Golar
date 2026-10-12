package main

import (
	"example/hello/db"
	"example/hello/startup"
)

func main() {
	db.StartDb()
	startup.Server()
}
