package router2

import (
	"log"
	"main/auth2"
	"main/db"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getBrowse(c *gin.Context) {
	courses, _ := db.GetAllPublicCoursesPreload()

	c.HTML(
		http.StatusOK,
		"browse_.html",
		gin.H{
			"User":     auth2.GetLoggedInUserLogError(c),
			"LoggedIn": auth2.IsLoggedInValid(c),
			"Messages": msg.GetMessages(c),
			"Courses":  courses,
		},
	)
}

func getCourse(c *gin.Context) {
	usernameParam := c.Params.ByName("username")
	courseParam := c.Params.ByName("course")
	sectionIDParam := c.Params.ByName("sectionID")

	course, err := db.GetUserCoursePreload(usernameParam, courseParam)

	if err != nil {
		log.Println("router/get.go ERROR getting course:", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// if authro then they can view private releases
	releases := []db.Release{}
	if auth2.GetLoggedInUserLogError(c).ID == course.UserID {
		releases, _ = db.GetAnyReleases(course.ID)
	} else {

		releases, _ = db.GetPublicReleases(course.ID)
	}

	c.HTML(
		http.StatusOK,
		"course_.html",
		gin.H{
			"SectionID": sectionIDParam,
			"User":      auth2.GetLoggedInUserLogError(c),
			"LoggedIn":  auth2.IsLoggedInValid(c),
			"Messages":  msg.GetMessages(c),
			"Releases":  releases,
			"Course":    course,
		},
	)
}
