package githubapi

import (
	"context"
	"log"
	"main/db"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

func GetVersionGithubTree(version *db.Version, c *gin.Context) (*github.Tree, error) {
	// get course owner
	course, err2 := db.GetCoursePreloadUser(version.CourseID)
	if err2 != nil {
		log.Println("api/github ERROR getting course in getGithubRepoCommitTree:", err2)
		return nil, err2
	}

	// get version's githubVersion
	githubVersion, err1 := version.GetGithubVersion()
	if err1 != nil {
		log.Println("api/github ERROR getting githubVersion in getGithubRepoCommitTree:", err1)
		return nil, err1
	}

	user, err3 := db.GetUser(course.UserID)
	if err3 != nil {
		log.Println("api/github ERROR getting user in getGithubRepoCommitTree:", err3)
		return nil, err3
	}

	// get owner's github connection
	connection, err4 := GetGithubConnection(user)
	if err4 != nil {
		log.Println("api/github ERROR getting user's github connection in getGithubRepoCommitTree:", err4)
		return nil, err4
	}

	ctx := context.Background()
	// get client
	client := NewClient(connection, ctx)

	githubUser, _, err5 := client.Users.Get(ctx, "")
	if err5 != nil {
		log.Println("api/github ERROR getting github user in getGithubRepoCommitTree:", err5)
		return nil, err5
	}

	repo, _, err6 := client.Repositories.GetByID(ctx, githubVersion.RepoID)
	if err6 != nil {
		log.Println("api/github ERROR getting repo by ID in getGithubRepoCommitTree:", err6)
		return nil, err6
	}

	// get folders from repo with info from githubVersion
	// use sha to get specific commit
	tree, _, err7 := client.Git.GetTree(ctx, *githubUser.Login, *repo.Name, githubVersion.SHA, true)
	if err7 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubRepoCommitTree:", err7)
		return nil, err7
	}

	return tree, nil
}

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

func GetGithubConnection(user *db.User) (*db.GithubConnection, error) {
	githubConnection := db.GithubConnection{}
	err := db.GormDB.Model(&db.GithubConnection{}).Where("user_id = ?", user.ID).First(&githubConnection).Error
	return &githubConnection, err
}

func GithubGetReposLogError(user *db.User) []*github.Repository {
	connection, err := GetGithubConnection(user)
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

func GithubGetRepoLogError(user *db.User, repoID uint64) *github.Repository {
	connection, err := GetGithubConnection(user)
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
