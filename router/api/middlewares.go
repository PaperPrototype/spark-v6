package api

import (
	"log"
	"main/db"
	"main/router/auth"
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

// uses course_id in url c.Query or c.PostForm to confirm if the currently logged in user is an author of the course
func mustBeCourseAuthor(c *gin.Context) {
	loggedInValid := auth.IsLoggedInValid(c)

	if !loggedInValid {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	courseID := c.PostForm("course_id")

	// if postform is empty try url query
	if courseID == "" {
		courseID = c.Query("course_id")
	}

	course, err1 := db.GetWithIDUserCoursePreloadUser(user.ID, courseID)
	if err1 != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if course.UserID != user.ID {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// if all is well then continue
	c.Next()
}
