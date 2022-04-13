package helpers

import (
	"strings"

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

func GetGithubClientID() string {
	if FileExists("./githubclientid") {
		data, err := os.ReadFile("./githubclientid")
		if err != nil {
			log.Println("config: error reading githubclientid file")
			panic(err)
		}

		return string(data)
	}

	env := os.Getenv("GITHUB_CLIENT_ID")

	if env == "" {
		panic(errors.New("empty env variable for GITHUB_CLIENT_ID"))
	}

	return env
}

func GetGithubClientSecret() string {
	if FileExists("./githubclientsecret") {
		data, err := os.ReadFile("./githubclientsecret")
		if err != nil {
			log.Println("config: error reading githubclientid file")
			panic(err)
		}

		return string(data)
	}

	env := os.Getenv("GITHUB_CLIENT_SECRET")

	if env == "" {
		panic(errors.New("empty env variable for GITHUB_CLIENT_SECRET"))
	}

	return env
}

const AllowedUsernameCharacters string = "abcdefghijklmnopqrstuvwxyz1234567890-_"

func IsAllowedUsername(username string) bool {
	for _, char := range username {
		if !strings.Contains(AllowedUsernameCharacters, string(char)) {
			return false
		}
	}
	return true
}

// get usernames after an @ out of markdown and return them
func GetUserMentions(textWithMentions string) []string {
	usernames := []string{}
	for i, c := range textWithMentions {
		if c == '@' {
			usernames = append(usernames, getUntilSpace(i+1, textWithMentions))
		}
	}

	return usernames
}

func getUntilSpace(index int, src string) string {
	name := []rune{}
	for index < len(src) {
		if src[index] == ' ' {
			break
		}
		name = append(name, rune(src[index]))
		index++
	}

	return string(name)
}
