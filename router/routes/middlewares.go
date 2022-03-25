package routes

import (
	"main/db"
	"main/router/session"

	"github.com/gin-gonic/gin"
)

func mustBeCourseEditor(c *gin.Context) {
	loggedInValid := session.IsLoggedInValid(c)

	if !loggedInValid {
		notFound(c)
		return
	}

	user, err := session.GetLoggedInUser(c)
	if err != nil {
		notFound(c)
		return
	}

	name := c.Params.ByName("course")

	course, err1 := db.GetCourse(name)
	if err1 != nil {
		notFound(c)
		return
	}

	if course.UserID != user.ID {
		notFound(c)
		return
	}

	// if all is well then continue
	c.Next()
}
