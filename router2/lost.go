package router2

import (
	"main/auth2"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getLost(c *gin.Context) {
	c.HTML(
		http.StatusNotFound,
		"lost_.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth2.GetLoggedInUserLogError(c),
			"LoggedIn": auth2.IsLoggedInValid(c),
			"Meta": meta{
				Title: "Sparker - 404d",
				Desc:  "404 page not found",
			},
		},
	)
}
