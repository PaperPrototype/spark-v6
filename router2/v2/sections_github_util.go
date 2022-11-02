package v2

import (
	"context"
	"log"
	"main/db"
	"main/githubapi"

	"github.com/google/go-github/github"
)

// get markdown from github
func getGithubMarkdown(author *db.User, release *db.Release, section *db.Section) (HTML string, Error string) {
	// get authors's github connection
	connection, err4 := githubapi.GetGithubConnection(author.ID)
	if err4 != nil {
		log.Println("v2/sections.go ERROR getting user's github connection in getSectionMarkdown:", err4)

		return "", "Error getting github connection for author"
	}

	ctx := context.Background()

	// get client
	client := githubapi.NewClient(connection, ctx)

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
	contentEncoded, _, _, err7 := client.Repositories.GetContents(ctx, *githubUser.Login, *repo.Name, section.GithubSection.Path, &github.RepositoryContentGetOptions{
		Ref: release.GithubRelease.SHA,
	})
	if err7 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubRepoCommitContent:", err7)
		return "", "Internal Server Error"
	}

	// decode content
	content, err8 := contentEncoded.GetContent()
	if err8 != nil {
		log.Println("api/github ERROR decoding", *contentEncoded.Encoding, "content in getGithubRepoCommitContent:", err8)
		return "", "Internal Server Error"
	}

	return string(content), ""
}
