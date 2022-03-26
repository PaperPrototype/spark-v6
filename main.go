package main

import (
	"main/conn"
	"main/db"
	"main/router"
)

func main() {
	router.Setup()

	db.Setup()

	conn.Setup()

	router.Run()

	conn.Close()
}
