package db

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (user *User) HasGithubConnection() bool {
	var count int64 = 0
	err := gormDB.Model(&GithubConnection{}).Where("user_id = ?", user.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/github ERROR getting github connection:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func (user *User) GetGithubConnectionLogError() *GithubConnection {
	githubConnection := GithubConnection{}
	err := gormDB.Model(&GithubConnection{}).Where("user_id = ?", user.ID).First(&githubConnection).Error
	if err != nil {
		log.Println("db/github ERROR getting github connection:", err)
	}

	return &githubConnection
}

func (githubConnection *GithubConnection) NewClient() (*github.Client, context.Context) {
	ctx := context.Background()

	// put token into oauth struct
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubConnection.AccessToken})

	// http client
	client := oauth2.NewClient(ctx, tokenSource)
	return github.NewClient(client), ctx
}

func (user *User) GetGithubConnection() (*GithubConnection, error) {
	githubConnection := GithubConnection{}
	err := gormDB.Model(&GithubConnection{}).Where("user_id = ?", user.ID).First(&githubConnection).Error
	return &githubConnection, err
}

func (user *User) GithubGetReposLogError() []*github.Repository {
	connection, err := user.GetGithubConnection()
	if err != nil {
		log.Println("db/github ERROR getting github connection in GetReposLogError:", err)
	}

	client, ctx := connection.NewClient()

	repos, _, err1 := client.Repositories.List(ctx, "", nil)
	if err1 != nil {
		log.Println("db/github ERROR getting repos in GetReposLogError", err1)
	}

	return repos
}
