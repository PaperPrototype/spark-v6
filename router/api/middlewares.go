package api

import (
	"log"
	"main/db"
	"main/router/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mustBeCourseEditor(c *gin.Context) {
	// must be logged in
	if !session.IsLoggedInValid(c) {
		// abort with internal error to hide the fact this course exists
		// for security of private courses
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := session.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api ERROR getting logged in user:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// middleware to authenticate user
	courseID := c.Params.ByName("courseID")
	course, err := db.GetCourseWithIDStr(courseID)
	if err != nil {
		log.Println("api ERROR getting course:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user.ID != course.ID {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Next()
}
