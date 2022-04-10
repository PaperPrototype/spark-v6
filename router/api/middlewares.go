package api

import (
	"main/db"
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
