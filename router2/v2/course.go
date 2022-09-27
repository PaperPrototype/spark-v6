package v2

import (
	"log"
	"main/auth2"
	"main/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCourseReleasesJSON(c *gin.Context) {
	courseID := c.Params.ByName("courseID")

	course, err := db.GetCourseWithIDPreloadUser(courseID)
	if err != nil {
		log.Println("v2/course.go ERROR getting course:", err)
		c.JSON(http.StatusInternalServerError, payload{
			Error: "Error getting course",
		})
		return
	}

	// if author then they can view private releases
	releases := []db.Release{}
	if auth2.GetLoggedInUserLogError(c).ID == course.UserID {
		releases, _ = db.GetAnyReleases(course.ID)
	} else {

		releases, _ = db.GetPublicReleases(course.ID)
	}

	c.JSON(http.StatusOK, payload{
		Payload: releases,
	})
}
