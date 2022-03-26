package helpers

import (
	"context"

	pgx "github.com/jackc/pgx/v4"

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

func GetDBConnTemp() *pgx.Conn {
	if FileExists("./dbconfig") {
		data, err := os.ReadFile("./dbconfig")
		if err != nil {
			log.Println("helpers ERROR reading dbconfig file:", err)
			panic(err)
		}

		conn, err1 := pgx.Connect(context.Background(), string(data))
		if err1 != nil {
			log.Println("helpers ERROR opening connection:", err1)
			panic(err)
		}

		return conn
	}

	env := os.Getenv("DATABASE_URL")

	if env == "" {
		panic(errors.New("empty env variable for DATABASE_URL"))
	}

	conn, err1 := pgx.Connect(context.Background(), string(env))
	if err1 != nil {
		log.Println("helpoers ERROR opening connection:", err1)
		panic(err1)
	}

	return conn
}
