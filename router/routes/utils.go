package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notFound(c *gin.Context) {
	c.Redirect(http.StatusFound, "/lost")
}
