package api

import (
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getGithubUserRepos(c *gin.Context) {
	user := auth.GetLoggedInUserLogError(c)

	c.JSON(
		http.StatusOK,
		user.GithubGetReposLogError(),
	)
}
