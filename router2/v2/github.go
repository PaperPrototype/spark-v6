package v2

import (
	"log"
	"main/auth2"
	"main/githubapi"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// group.GET("/v2/user/github/repos", mustBeLoggedIn)
func getUserGithubReposJSON(c *gin.Context) {
	user := auth2.GetLoggedInUserLogError(c)

	c.JSON(
		http.StatusOK,
		payload{
			Payload: githubapi.GithubGetReposLogError(user),
		},
	)
}

// group.GET("/v2/user/github/repo/:repoID/branches", mustBeLoggedIn)
func getUserGithubRepoBranchesJSON(c *gin.Context) {
	repoID := c.Params.ByName("repoID")

	user := auth2.GetLoggedInUserLogError(c)
	connection, err := githubapi.GetGithubConnection(user)
	if err != nil {
		log.Println("api/github ERROR getting github connection in getGithubRepoBranches:", err)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github connection",
			},
		)
		return
	}

	repoIDNum, err1 := strconv.ParseInt(repoID, 10, 64)
	if err1 != nil {
		log.Println("api/github ERROR parsing repoID in getGithubRepoBranches:", err1)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Parse error",
			},
		)
		return
	}

	branches, err2 := githubapi.GetRepoByIDBranches(connection, repoIDNum)
	if err2 != nil {
		log.Println("api/github ERROR getting repo in getGithubRepoBranches:", err2)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting repository branches",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		payload{
			Error:   "",
			Payload: branches,
		},
	)
}

// group.GET("/v2/user/github/repo/:repoID/branch/:branch/commits", mustBeLoggedIn)
func getUserGithubRepoBranchCommitsJSON(c *gin.Context) {
	repoID := c.Params.ByName("repoID")
	branch := c.Params.ByName("branch")

	user := auth2.GetLoggedInUserLogError(c)

	connection, err := githubapi.GetGithubConnection(user)
	if err != nil {
		log.Println("apit/github ERROR getting github connection in getGithubRepoBranchCommits:", err)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github connection",
			},
		)
		return
	}

	repoIDNum, err1 := strconv.ParseInt(repoID, 10, 64)
	if err1 != nil {
		log.Println("api/github ERROR parsing repoID in getGithubRepoBranchCommits:", err1)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Parse error",
			},
		)
		return
	}

	commits, err2 := githubapi.GetCommitsByRepoIDBranch(connection, repoIDNum, branch)
	if err2 != nil {
		log.Println("api/github ERROR getting commits by RepoID and Branch in getGithubRepoBranchCommits:", err2)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting commits",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		payload{
			Payload: commits,
		},
	)
}
