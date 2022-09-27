package middlewares

import (
	"main/auth2"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MustBeLoggedIn(c *gin.Context) {
	if !auth2.IsLoggedInValid(c) {
		msg.SendMessage(c, "We kinda need you to be logged in to access that page...")
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.Next()
}
