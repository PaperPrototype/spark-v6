package routes

import (
	"log"
	"main/db"
	"main/msg"
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mustBeCourseEditor(c *gin.Context) {
	loggedInValid := auth.IsLoggedInValid(c)

	if !loggedInValid {
		notFound(c)
		return
	}

	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		notFound(c)
		return
	}

	name := c.Params.ByName("course")
	username := c.Params.ByName("username")
	course, err1 := db.GetUserCoursePreloadUser(username, name)
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

// OVERVIEW
// if not owner of course
// 		if course release not public
// 				return
// 		if course release not free
//			if not payed
// 				return and pay
// coninue
func MustHaveAccessToCourseRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	versionID := c.Params.ByName("versionID")

	course, err2 := db.GetUserCoursePreloadUser(username, courseName)
	if err2 != nil {
		log.Println("routes/MustHaveAccessToCourseRelease ERROR getting course:", err2)
		notFound(c)
		return
	}

	redirect := false
	version, err := course.GetVersion(versionID)
	if err != nil {
		log.Println("routes/MustHaveAccessToCourseRelease ERROR getting version:", err)
		msg.SendMessage(c, "That course upload may have been deleted.")
		redirect = true
	}

	if redirect {
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	redirectRelease := false
	release, err1 := db.GetAllRelease(version.ReleaseID)
	if err1 != nil {
		log.Println("routes/MustHaveAccessToCourseRelease ERROR getting release:", err1)
		msg.SendMessage(c, "No course releases available.")
		redirectRelease = true
	}

	if redirectRelease {
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	user := auth.GetLoggedInUserLogError(c)

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
			if !db.UserCanAccessCourseRelease(user.ID, version) {
				msg.SendMessage(c, "You kinda have to pay to access that")
				c.Redirect(http.StatusFound, "/"+username+"/"+course.Name)
				return
			}
		}
	}

	// coninue
	c.Next()
}
