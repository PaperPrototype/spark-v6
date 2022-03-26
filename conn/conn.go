package conn

import (
	"context"
	"errors"
	"log"
	"main/helpers"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func Setup() {
	if helpers.FileExists("./dbconfig") {
		data, err := os.ReadFile("./dbconfig")
		if err != nil {
			log.Println("helpers ERROR reading dbconfig file:", err)
			panic(err)
		}

		poolTmp, err1 := pgxpool.Connect(context.Background(), string(data))
		if err1 != nil {
			log.Println("helpers ERROR opening connection:", err1)
			panic(err)
		}

		pool = poolTmp
		return
	}

	env := os.Getenv("DATABASE_URL")

	if env == "" {
		panic(errors.New("empty env variable for DATABASE_URL"))
	}

	poolTmp, err1 := pgxpool.Connect(context.Background(), string(env))
	if err1 != nil {
		log.Println("helpers ERROR opening connection:", err1)
		panic(err1)
	}

	pool = poolTmp
}

func Close() {
	pool.Close()
}

func GetConn() (*pgxpool.Conn, error) {
	return pool.Acquire(context.Background())
}
