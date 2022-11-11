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

// get markdown from github
func GetGithubMarkdown(author *db.User, release *db.Release, path string) (HTML string, Error string) {
	// get authors's github connection
	connection, err4 := GetGithubConnection(author.ID)
	if err4 != nil {
		log.Println("v2/sections.go ERROR getting user's github connection in getSectionMarkdown:", err4)

		return "", "Error getting github connection for author"
	}

	ctx := context.Background()

	// get client
	client := NewClient(connection, ctx)

	githubUser, _, err5 := client.Users.Get(ctx, "")
	if err5 != nil {
		log.Println("v2/sections.go ERROR getting github user in getSectionMarkdown:", err5)

		return "", "Error getting github user"
	}

	repo, _, err6 := client.Repositories.GetByID(ctx, int64(release.GithubRelease.RepoID))
	if err6 != nil {
		log.Println("v2/sections.go ERROR getting repo by ID in getSectionMarkdown:", err6)

		return "", "Error getting github user"
	}

	// get folders from repo with info from githubVersion
	// use sha to get specific commit
	contentEncoded, _, _, err7 := client.Repositories.GetContents(ctx, *githubUser.Login, *repo.Name, path, &github.RepositoryContentGetOptions{
		Ref: release.GithubRelease.SHA,
	})
	if err7 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubRepoCommitContent:", err7)
		return "", "Error getting markdown at path. Did you rename a file or folder in the github repository?"
	}

	// decode content
	content, err8 := contentEncoded.GetContent()
	if err8 != nil {
		log.Println("api/github ERROR decoding", *contentEncoded.Encoding, "content in getGithubRepoCommitContent:", err8)
		return "", "Internal Server Error Decoding Contents"
	}

	return string(content), ""
}
