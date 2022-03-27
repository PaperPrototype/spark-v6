package routes

import (
	"log"
	"main/db"
	"main/msg"
	"main/router/session"
	"net/http"

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

func mustBeLoggedIn(c *gin.Context) {
	if !session.IsLoggedInValid(c) {
		msg.SendMessage(c, "We kinda need you to be logged in to access that page...")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Next()
}

// OVERVIEW
// if not owner of course
// 		if course release not public
// 				return
// 		if course release not free
//			if not payed
// 				return and pay
// coninue
func MustHaveAccessToCourseRelease(c *gin.Context) {
	courseName := c.Params.ByName("course")
	versionID := c.Params.ByName("versionID")

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("routes/MustHaveAccessToCourseRelease ERROR getting version:", err)
		notFound(c)
		return
	}

	release, err1 := db.GetRelease(version.ReleaseID)
	if err1 != nil {
		log.Println("routes/MustHaveAccessToCourseRelease ERROR getting release:", err1)
		notFound(c)
		return
	}

	course, err2 := db.GetCourse(courseName)
	if err2 != nil {
		log.Println("routes/MustHaveAccessToCourseRelease ERROR getting course:", err2)
		notFound(c)
		return
	}

	user := session.GetLoggedInUserHideError(c)

	// if not owner of course
	if course.UserID != user.ID {
		// if course release not public
		if !release.Public {
			notFound(c)
			return
		}

		// if course release not free
		if release.Price != 0 {
			// if not payed
			if !db.UserHasPurchasedCourseRelease(user.ID, release.ID) {
				msg.SendMessage(c, "You kinda have to pay to access that")
				c.Redirect(http.StatusFound, "/"+course.Name)
				return
			}
		}
	}

	// coninue
	c.Next()
}
