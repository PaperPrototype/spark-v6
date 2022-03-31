package api

import (
	"log"
	"main/db"
	session "main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Un-used middleware to authenticate user
func ustBeCourseEditor(c *gin.Context) {
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

	courseID := c.Params.ByName("courseID")
	course, err := db.GetCourseWithIDPreloadUser(courseID)
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
