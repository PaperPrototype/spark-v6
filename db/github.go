package db

import (
	"context"
	"log"
	"main/githubapi"

	"github.com/google/go-github/github"
)

func (user *User) HasGithubConnection() bool {
	var count int64 = 0
	err := gormDB.Model(&githubapi.GithubConnection{}).Where("user_id = ?", user.ID).Count(&count).Error

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

func (user *User) GetGithubConnectionLogError() *githubapi.GithubConnection {
	githubConnection := githubapi.GithubConnection{}
	err := gormDB.Model(&githubapi.GithubConnection{}).Where("user_id = ?", user.ID).First(&githubConnection).Error
	if err != nil {
		log.Println("db/github ERROR getting github connection:", err)
	}

	return &githubConnection
}

func (user *User) GetGithubConnection() (*githubapi.GithubConnection, error) {
	githubConnection := githubapi.GithubConnection{}
	err := gormDB.Model(&githubapi.GithubConnection{}).Where("user_id = ?", user.ID).First(&githubConnection).Error
	return &githubConnection, err
}

func (user *User) GithubGetReposLogError() []*github.Repository {
	connection, err := user.GetGithubConnection()
	if err != nil {
		log.Println("db/github ERROR getting github connection in GetReposLogError:", err)
	}

	ctx := context.Background()
	client := connection.NewClient(ctx)

	repos, _, err1 := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
		Visibility: "all",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	if err1 != nil {
		log.Println("db/github ERROR getting repos in GetReposLogError", err1)
	}

	return repos
}

func GetGithubReleaseWithIDStr(releaseID string) (*GithubRelease, error) {
	githubRelease := GithubRelease{}
	err := gormDB.Model(&GithubRelease{}).Where("release_id = ?", releaseID).First(&githubRelease).Error
	if err != nil {
		return &githubRelease, err
	}
	return &githubRelease, nil
}
