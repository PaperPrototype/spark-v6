package githubapi

import (
	"context"
	"log"
	"main/db"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func NewClient(githubConnection *db.GithubConnection, context context.Context) *github.Client {
	// put token into oauth struct
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubConnection.AccessToken})

	// http client
	client := oauth2.NewClient(context, tokenSource)
	return github.NewClient(client)
}

func GetRepoByIDBranches(githubConnection *db.GithubConnection, repoID int64) ([]*github.Branch, error) {
	ctx := context.Background()
	client := NewClient(githubConnection, ctx)

	repo, _, err1 := client.Repositories.GetByID(ctx, repoID)
	if err1 != nil {
		log.Println("githubapi ERROR getting repo in GetRepoByIDBranches:", err1)
		return []*github.Branch{}, err1
	}

	branches, _, err2 := client.Repositories.ListBranches(ctx, *repo.GetOwner().Login, *repo.Name, &github.ListOptions{Page: 1, PerPage: 100})
	if err2 != nil {
		log.Println("githubapi ERROR listing branches in GetRepoByIDBranches:", err2)
		return branches, err2
	}

	return branches, nil
}

func GetRepoBranch(githubConnection *db.GithubConnection, repo string, branch string) (*github.Branch, error) {
	ctx := context.Background()
	client := NewClient(githubConnection, ctx)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println("githubapi ERROR getting user in GetRepoBranch:", err)
		return nil, err
	}

	branchRepo, _, err1 := client.Repositories.GetBranch(ctx, *user.Login, repo, branch)
	if err1 != nil {
		log.Println("githubapi ERROR getting repo in GetRepoByIDBranches:", err1)
	}

	return branchRepo, err1
}

func GetCommitsByRepoIDBranch(githubConnection *db.GithubConnection, repoID int64, branch string) ([]*github.RepositoryCommit, error) {
	ctx := context.Background()
	client := NewClient(githubConnection, ctx)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println("githubapi ERROR getting user in GetRepoByIDBranch:", err)
		return nil, err
	}

	repo, _, err2 := client.Repositories.GetByID(ctx, repoID)
	if err2 != nil {
		log.Println("githubapi ERROR getting repo in GetRepoByIDBranch:", err2)
		return nil, err2
	}

	branches, _, err1 := client.Repositories.ListCommits(ctx, *user.Login, *repo.Name, &github.CommitsListOptions{
		SHA: branch,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	if err1 != nil {
		log.Println("githubapi ERROR getting repo in GetRepoByIDBranch:", err1)
	}

	return branches, err1
}

func GetRepoCommit(githubConnection *db.GithubConnection, repo string, sha string) (*github.RepositoryCommit, error) {
	ctx := context.Background()
	client := NewClient(githubConnection, ctx)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println("githubapi ERROR getting user in GetRepoCommit:", err)
		return nil, err
	}

	repoCommit, _, err1 := client.Repositories.GetCommit(ctx, *user.Login, repo, sha)
	if err1 != nil {
		log.Println("githubapi ERROR getting repo commit in GetRepoCommit:", err1)
	}

	return repoCommit, err1
}

func GetRepoCommits(githubConnection *db.GithubConnection, repo string, branch string) ([]*github.RepositoryCommit, error) {
	ctx := context.Background()
	client := NewClient(githubConnection, ctx)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println("githubapi ERROR getting user in GetRepoCommits:", err)
		return nil, err
	}

	commits, _, err1 := client.Repositories.ListCommits(ctx, *user.Login, repo, &github.CommitsListOptions{
		SHA: branch,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	})
	if err1 != nil {
		log.Println("githubapi ERROR getting repo commits in GetRepoCommits:", err1)
	}

	return commits, err1
}
