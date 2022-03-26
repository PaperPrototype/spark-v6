package helpers

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"

	"errors"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const passwordHashCostAndStrength int = 8

func HashPassword(pass string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), passwordHashCostAndStrength)
	return string(hashPass), err
}

func GetDatabaseURL() string {
	if FileExists("./dbconfig") {
		data, err := os.ReadFile("./dbconfig")
		if err != nil {
			log.Println("config: error reading dbconfig file")
			panic(err)
		}

		return string(data)
	}

	env := os.Getenv("DATABASE_URL")

	if env == "" {
		panic(errors.New("empty env variable for DATABASE_URL"))
	}

	return env
}

func FileExists(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		log.Println("helper: NO DBCONFIG FILE!!!", err)
		return false
	}

	return true
}

// MAKE SURE TO CLOSE THE CONNECTION ONCE YOU ARE DONE!
func GetDBConnTemp() *sql.Conn {
	if FileExists("./dbconfig") {
		data, err := os.ReadFile("./dbconfig")
		if err != nil {
			log.Println("helpers ERROR reading dbconfig file:", err)
			panic(err)
		}

		db, err1 := sql.Open("pgx", string(data))
		if err1 != nil {
			log.Println("helpers ERROR opening connection:", err1)
			panic(err)
		}

		conn, err2 := db.Conn(context.Background())
		if err2 != nil {
			panic(err2)
		}

		return conn
	}

	env := os.Getenv("DATABASE_URL")

	if env == "" {
		panic(errors.New("empty env variable for DATABASE_URL"))
	}

	db, err1 := sql.Open("pgx", string(env))
	if err1 != nil {
		log.Println("helpers ERROR opening connection:", err1)
		panic(err1)
	}

	conn, err2 := db.Conn(context.Background())
	if err2 != nil {
		panic(err2)
	}

	return conn
}
