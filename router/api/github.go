package api

import (
	"context"
	"log"
	"main/db"
	"main/markdown"
	"main/router/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

func getGithubUserRepos(c *gin.Context) {
	user := auth.GetLoggedInUserLogError(c)

	c.JSON(
		http.StatusOK,
		user.GithubGetReposLogError(),
	)
}

func getGithubRepoBranches(c *gin.Context) {
	repoID := c.Params.ByName("repoID")

	user := auth.GetLoggedInUserLogError(c)
	connection, err := user.GetGithubConnection()
	if err != nil {
		log.Println("api/github ERROR getting github connection in getGithubRepoBranches:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	repoIDNum, err1 := strconv.ParseInt(repoID, 10, 64)
	if err1 != nil {
		log.Println("api/github ERROR parsing repoID in getGithubRepoBranches:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	branches, err2 := connection.GetRepoByIDBranches(repoIDNum)
	if err2 != nil {
		log.Println("api/github ERROR getting repo in getGithubRepoBranches:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		branches,
	)
}

func getGithubRepoBranchCommits(c *gin.Context) {
	repoID := c.Params.ByName("repoID")
	branch := c.Params.ByName("branch")

	user := auth.GetLoggedInUserLogError(c)

	connection, err := user.GetGithubConnection()
	if err != nil {
		log.Println("apit/github ERROR getting github connection in getGithubRepoBranchCommits:", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	repoIDNum, err1 := strconv.ParseInt(repoID, 10, 64)
	if err1 != nil {
		log.Println("api/github ERROR parsing repoID in getGithubRepoBranchCommits:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	commits, err2 := connection.GetCommitsByRepoIDBranch(repoIDNum, branch)
	if err2 != nil {
		log.Println("api/github ERROR getting commits by RepoID and Branch in getGithubRepoBranchCommits:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		commits,
	)
}

func getGithubRepoCommitTree(c *gin.Context) {
	/*
		URL params
		/api/github/version/:versionID/trees/:tree_sha
	*/
	versionID := c.Params.ByName("versionID")

	if versionID == "" {
		log.Println("api/github ERROR versionID is empty.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get version
	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api/github ERROR getting version in getGithubRepoCommitTree:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get course owner
	course, err2 := db.GetCourse(version.CourseID)
	if err2 != nil {
		log.Println("api/github ERROR getting course in getGithubRepoCommitTree:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get version's githubVersion
	githubVersion, err1 := version.GetGithubVersion()
	if err1 != nil {
		log.Println("api/github ERROR getting githubVersion in getGithubRepoCommitTree:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err3 := db.GetUser(course.UserID)
	if err3 != nil {
		log.Println("api/github ERROR getting user in getGithubRepoCommitTree:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get owner's github connection
	connection, err4 := user.GetGithubConnection()
	if err4 != nil {
		log.Println("api/github ERROR getting user's github connection in getGithubRepoCommitTree:", err4)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	// get client
	client := connection.NewClient(ctx)

	githubUser, _, err5 := client.Users.Get(ctx, "")
	if err5 != nil {
		log.Println("api/github ERROR getting github user in getGithubRepoCommitTree:", err5)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	repo, _, err6 := client.Repositories.GetByID(ctx, githubVersion.RepoID)
	if err6 != nil {
		log.Println("api/github ERROR getting repo by ID in getGithubRepoCommitTree:", err6)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get folders from repo with info from githubVersion
	// use sha to get specific commit
	tree, _, err7 := client.Git.GetTree(ctx, *githubUser.Login, *repo.Name, githubVersion.SHA, true)
	if err7 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubRepoCommitTree:", err7)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		tree,
	)
}

func getGithubRepoCommitContent(c *gin.Context) {
	/*
		/api/github/version/:versionID/content/:commit_sha/*path
	*/

	versionID := c.Params.ByName("versionID")
	commitSHA := c.Params.ByName("commit_sha")
	path := c.Params.ByName("path")

	if versionID == "" {
		log.Println("api/github ERROR versionID is empty.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get version
	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api/github ERROR getting version in getGithubRepoCommitContent:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get course owner
	course, err2 := db.GetCourse(version.CourseID)
	if err2 != nil {
		log.Println("api/github ERROR getting course in getGithubRepoCommitContent:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get version's githubVersion
	githubVersion, err1 := version.GetGithubVersion()
	if err1 != nil {
		log.Println("api/github ERROR getting githubVersion in getGithubRepoCommitContent:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err3 := db.GetUser(course.UserID)
	if err3 != nil {
		log.Println("api/github ERROR getting user in getGithubRepoCommitContent:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get owner's github connection
	connection, err4 := user.GetGithubConnection()
	if err4 != nil {
		log.Println("api/github ERROR getting user's github connection in getGithubRepoCommitContent:", err4)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	// get client
	client := connection.NewClient(ctx)

	githubUser, _, err5 := client.Users.Get(ctx, "")
	if err5 != nil {
		log.Println("api/github ERROR getting github user in getGithubRepoCommitContent:", err5)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	repo, _, err6 := client.Repositories.GetByID(ctx, githubVersion.RepoID)
	if err6 != nil {
		log.Println("api/github ERROR getting repo by ID in getGithubRepoCommitContent:", err6)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// get folders from repo with info from githubVersion
	// use sha to get specific commit
	contentEncoded, _, _, err7 := client.Repositories.GetContents(ctx, *githubUser.Login, *repo.Name, path, &github.RepositoryContentGetOptions{
		Ref: commitSHA,
	})
	if err7 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubRepoCommitContent:", err7)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// decode content
	content, err8 := contentEncoded.GetContent()
	if err8 != nil {
		log.Println("api/github ERROR decoding", *contentEncoded.Encoding, "content in getGithubRepoCommitContent:", err8)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	html, err9 := markdown.Convert([]byte(content))
	if err9 != nil {
		log.Println("api/github ERROR converting content ot markdown in getGithubRepoCommitContent:", err9)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	name := *contentEncoded.Name

	if len(name) != 0 {
		name = name[0 : len(*contentEncoded.Name)-3]
	}

	c.JSON(
		http.StatusOK,
		struct {
			Name     string
			Markdown string
		}{
			Name:     name,
			Markdown: html.String(),
		},
	)
}
