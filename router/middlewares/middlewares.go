package middlewares

import (
	"main/msg"
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MustBeLoggedIn(c *gin.Context) {
	if !auth.IsLoggedInValid(c) {
		msg.SendMessage(c, "We kinda need you to be logged in to access that page...")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Next()
}
