package router2

import (
	"log"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getBuySuccess(c *gin.Context) {
	username := c.Params.ByName("username")
	course := c.Params.ByName("course")

	log.Println("successfully purchased course (Your purchase may take a few minutes to complete).")
	msg.SendMessage(c, "Successfully purchased course")
	c.Redirect(http.StatusFound, "/"+username+"/"+course)
}

func getBuyCancel(c *gin.Context) {
	username := c.Params.ByName("username")
	course := c.Params.ByName("course")

	log.Println("canceled course purchase")
	msg.SendMessage(c, "Purchase canceled")
	c.Redirect(http.StatusFound, "/"+username+"/"+course)
}
