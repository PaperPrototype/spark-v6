package v2

import (
	"fmt"
	"log"
	"main/auth2"
	"main/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mustBeSameUserID(c *gin.Context) {
	user, err := auth2.GetLoggedInUser(c)
	if err != nil {
		c.JSON(http.StatusOK, payload{
			Error: "Error getting logged in user. You must be logged in.",
		})
		return
	}

	userID := c.Params.ByName("userID")
	if userID != fmt.Sprint(user.ID) {
		c.JSON(http.StatusOK, payload{
			Error: "User ID doesn't match.",
		})
		return
	}

	c.Next()
}

func mustBeLoggedIn(c *gin.Context) {
	if !auth2.IsLoggedInValid(c) {
		c.JSON(http.StatusOK, payload{
			Error: "Error getting logged in user. You must be logged in.",
		})
		return
	}

	c.Next()
}

func mustBeAuthorSectionID(c *gin.Context) {
	if !auth2.IsLoggedInValid(c) {
		c.JSON(http.StatusOK, payload{
			Error: "Error getting logged in user. You must be logged in.",
		})
		return
	}

	user := auth2.GetLoggedInUserLogError(c)

	sectionID := c.Params.ByName("sectionID")

	section, err := db.GetSection(sectionID)
	if err != nil {
		c.JSON(http.StatusOK, payload{
			Error: "Error getting section",
		})
		return
	}

	release, err1 := db.GetAnyRelease(section.ReleaseID)
	if err1 != nil {
		log.Println("v2/middlewares.go ERROR getting any release in mustBeAuthorSectionID:", err1)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting release for that section",
		})
		return
	}

	course, err2 := db.GetCourse(release.CourseID)
	if err2 != nil {
		log.Println("v2/middlewares.go ERROR getting course in mustBeAuthorSectionID:", err2)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting course for that section",
		})
		return
	}

	if course.UserID != user.ID {
		c.JSON(http.StatusOK, payload{
			Error: "You must be the author of the course to edit a section",
		})
		return
	}

	c.Next()
}

func mustBeAuthorReleaseID(c *gin.Context) {
	if !auth2.IsLoggedInValid(c) {
		c.JSON(http.StatusOK, payload{
			Error: "Error getting logged in user. You must be logged in.",
		})
		return
	}

	user := auth2.GetLoggedInUserLogError(c)

	releaseID := c.Params.ByName("releaseID")
	release, err1 := db.GetAnyRelease(releaseID)
	if err1 != nil {
		log.Println("v2/middlewares.go ERROR getting any release in mustBeAuthorReleaseID:", err1)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting release",
		})
		return
	}

	course, err2 := db.GetCourse(release.CourseID)
	if err2 != nil {
		log.Println("v2/middlewares.go ERROR getting course in mustBeAuthorReleaseID:", err2)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting course",
		})
		return
	}

	if course.UserID != user.ID {
		c.JSON(http.StatusOK, payload{
			Error: "You must be the author of the course to edit this",
		})
		return
	}

	c.Next()
}

func mustBeAuthorCourseID(c *gin.Context) {
	if !auth2.IsLoggedInValid(c) {
		c.JSON(http.StatusOK, payload{
			Error: "Error getting logged in user. You must be logged in.",
		})
		return
	}

	user := auth2.GetLoggedInUserLogError(c)

	courseID := c.Params.ByName("courseID")

	course, err2 := db.GetCoursePreloadUser(courseID)
	if err2 != nil {
		log.Println("v2/middlewares.go ERROR getting course in mustBeAuthorCourseID:", err2)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting course",
		})
		return
	}

	if course.UserID != user.ID {
		c.JSON(http.StatusOK, payload{
			Error: "You must be the author of the course to edit this",
		})
		return
	}

	c.Next()
}
