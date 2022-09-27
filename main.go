package main

import (
	"main/conn"
	"main/db"
	"main/mailer"
	"main/payments"
	"main/router"
	"main/router2"
	"os"
)

func main() {
	args := os.Args

	if len(args) > 1 {
		// new router

		// db connection pool used for course uploading
		// and for performant db when ORM is too heavy
		conn.Setup()

		mailer.Setup()

		payments.Setup()

		db.Setup()

		router2.Run()

		conn.Close()
	} else {
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
}
