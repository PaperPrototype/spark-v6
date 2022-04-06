package api

import (
	"log"
	"main/router/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
