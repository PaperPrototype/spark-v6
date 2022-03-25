package main

import (
	"main/db"
	"main/router"
)

func main() {
	router.Setup()

	db.Setup()

	router.Run()
}
