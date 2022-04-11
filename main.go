package main

import (
	"main/conn"
	"main/db"
	"main/mailer"
	"main/payments"
	"main/router"
)

func main() {
	// db connection pool used for course uploading
	// and for performant db when ORM is too heavy
	conn.Setup()

	mailer.Setup()

	payments.Setup()

	router.Setup()

	db.Setup()

	router.Run()

	conn.Close()
}
