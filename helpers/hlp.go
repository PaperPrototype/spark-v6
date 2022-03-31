package helpers

import (
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
		return false
	}

	return true
}

func GetStripeKey() string {
	if FileExists("./stripeconfig") {
		data, err := os.ReadFile("./stripeconfig")
		if err != nil {
			log.Println("config: error reading stripeconfig file")
			panic(err)
		}

		return string(data)
	}

	env := os.Getenv("STRIPE_KEY")

	if env == "" {
		panic(errors.New("empty env variable for STRIPE_KEY"))
	}

	return env
}

func GetHost() string {
	env := os.Getenv("HOST_URL")
	if env == "" {
		env = "http://localhost:8080"
	}

	return env
}

func GetSendgridKey() string {
	if FileExists("./sendgridconfig") {
		data, err := os.ReadFile("./sendgridconfig")
		if err != nil {
			log.Println("config: error reading stripeconfig file")
			panic(err)
		}

		return string(data)
	}

	env := os.Getenv("SENDGRID_API_KEY")

	if env == "" {
		panic(errors.New("empty env variable for SENDGRID_API_KEY"))
	}

	return env
}
