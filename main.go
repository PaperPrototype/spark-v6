package main

import (
	"main/conn"
	"main/db"
	"main/mailer"
	"main/payments"
	"main/router2"
)

func main() {
	// new router

	// db connection pool used for course uploading
	// and for performant db when ORM is too heavy
	conn.Setup()

	mailer.Setup()

	payments.Setup()

	db.Setup()

	router2.Run()

	conn.Close()

}
