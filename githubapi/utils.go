package githubapi

import (
	"context"
	"log"
	"main/db"

	"github.com/google/go-github/github"
)

func GetGithubConnectionLogError(user *db.User) *db.GithubConnection {
	githubConnection := db.GithubConnection{}
	err := db.GormDB.Model(&db.GithubConnection{}).Where("user_id = ?", user.ID).First(&githubConnection).Error
	if err != nil {
		log.Println("db/github ERROR getting github connection:", err)
	}

	return &githubConnection
}

func UpdateGithubConnection(user *db.User, AccessToken string, TokenType string) error {
	return db.GormDB.Model(&db.GithubConnection{}).Where("user_id = ?", user.ID).Update("access_token", AccessToken).Update("token_type", TokenType).Error
}

func GetGithubConnection(userID uint64) (*db.GithubConnection, error) {
	githubConnection := db.GithubConnection{}
	err := db.GormDB.Model(&db.GithubConnection{}).Where("user_id = ?", userID).First(&githubConnection).Error
	return &githubConnection, err
}

func GithubGetReposLogError(userID uint64) []*github.Repository {
	connection, err := GetGithubConnection(userID)
	if err != nil {
		log.Println("db/github ERROR getting github connection in GetReposLogError:", err)
	}

	ctx := context.Background()
	client := NewClient(connection, ctx)

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

func GithubGetRepoLogError(userID uint64, repoID uint64) *github.Repository {
	connection, err := GetGithubConnection(userID)
	if err != nil {
		log.Println("db/github ERROR getting github connection in GetReposLogError:", err)
	}

	ctx := context.Background()
	client := NewClient(connection, ctx)

	repo, _, err1 := client.Repositories.GetByID(ctx, int64(repoID))
	if err1 != nil {
		log.Println("db/github ERROR getting repos in GetReposLogError", err1)
	}

	return repo
}

func GetGithubReleaseWithIDStr(user *db.User, releaseID string) (*db.GithubRelease, error) {
	githubRelease := db.GithubRelease{}
	err := db.GormDB.Model(&db.GithubRelease{}).Where("release_id = ?", releaseID).First(&githubRelease).Error
	if err != nil {
		return &githubRelease, err
	}
	return &githubRelease, nil
}
