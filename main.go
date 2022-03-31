package main

import (
	"main/conn"
	"main/db"
	"main/mailer"
	"main/payments"
	"main/router"
)

func main() {
	// mailer.Setup()
	mailer.Setup()

	payments.Setup()

	router.Setup()

	db.Setup()

	conn.Setup()

	router.Run()

	conn.Close()
}
