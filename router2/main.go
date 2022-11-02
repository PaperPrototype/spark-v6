package router2

import (
	"log"
	"main/auth2"
	"main/db"
	"main/msg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getBrowse(c *gin.Context) {
	courses, _ := db.GetAllPublicCoursesPreload()

	ownedCourses := []db.Ownership{}
	user, err := auth2.GetLoggedInUser(c)
	if err == nil {
		ownedCourses, _ = db.GetOwnershipsPreloadCourses(user.ID)
	}

	authoredCourses := []db.Course{}
	if err == nil {
		authoredCourses, _ = user.GetPublicAndPrivateAuthoredCourses()
	}

	c.HTML(
		http.StatusOK,
		"browse_.html",
		gin.H{
			"AuthoredCourses": authoredCourses,
			"OwnedCourses":    ownedCourses,
			"User":            user,
			"LoggedIn":        auth2.IsLoggedInValid(c),
			"Messages":        msg.GetMessages(c),
			"Courses":         courses,
			"Meta": meta{
				Title:    "Sparker - Browse",
				Desc:     "learn coding to build ideas",
				ImageURL: "/resources2/images/sparker_code_hl_banner.png",
			},
		},
	)
}

func getCourse(c *gin.Context) {
	usernameParam := c.Params.ByName("username")
	courseParam := c.Params.ByName("course")
	sectionIDParam := c.Params.ByName("sectionID")
	_ = c.Params.ByName("releaseID")

	course, err := db.GetUserCoursePreload(usernameParam, strings.ToLower(courseParam))

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

	owned := false
	user, err1 := auth2.GetLoggedInUser(c)
	if err1 != nil {
		if user.OwnsRelease(course.Release.ID) {
			owned = true
		}
	}

	sectionID := uint64(0)
	if sectionIDParam != "" {
		sectionID, _ = strconv.ParseUint(sectionIDParam, 10, 64)
	}

	c.HTML(
		http.StatusOK,
		"course_.html",
		gin.H{
			"Owned":     owned,
			"SectionID": sectionID,
			"User":      auth2.GetLoggedInUserLogError(c),
			"LoggedIn":  auth2.IsLoggedInValid(c),
			"Messages":  msg.GetMessages(c),
			"Releases":  releases,
			"Course":    course,
			"Meta": meta{
				Title:    course.Title,
				Desc:     course.Subtitle,
				ImageURL: course.Release.ImageURL,
			},
		},
	)
}
