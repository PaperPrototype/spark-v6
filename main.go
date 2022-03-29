package main

import (
	"main/conn"
	"main/db"
	"main/payments"
	"main/router"
)

func main() {
	payments.Setup()

	router.Setup()

	db.Setup()

	conn.Setup()

	router.Run()

	conn.Close()
}
