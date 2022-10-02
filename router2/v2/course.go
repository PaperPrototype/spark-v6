package v2

import (
	"log"
	"main/auth2"
	"main/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postCourseFORM(c *gin.Context) {
	courseID := c.Params.ByName("courseID")

	user := auth2.GetLoggedInUserLogError(c)
	title := c.PostForm("title")
	name := c.PostForm("name")
	public := c.PostForm("public") == "true"
	subtitle := c.PostForm("subtitle")

	available, err := db.UserCourseNameAvailableNotSelf(user.Username, name, courseID)
	if !available {
		c.JSON(http.StatusOK, payload{
			Error: "That course url name is taken",
		})
		return
	}

	if err != nil {
		log.Println("v2/course.go ERROR checking if course name is available in postCourseFORM:", err)
		c.JSON(http.StatusOK, payload{
			Error: "Error checking if course name available",
		})
		return
	}

	releasesCount := db.CountPublicCourseReleasesLogError(courseID)

	// if public and no releases!
	if public && releasesCount == 0 {
		c.JSON(http.StatusOK, payload{
			Error: "You must have at least 1 public release before you can make a course public",
		})
		return
	}

	err2 := db.UpdateCourse(courseID, title, name, subtitle, public)
	if err2 != nil {
		log.Println("v2/course.go ERROR updating course in postCourseFORM:", err2)
		c.JSON(http.StatusOK, payload{
			Error: "Error updating course",
		})
		return
	}

	c.JSON(http.StatusOK, payload{
		Payload: db.Course{
			UserID:   user.ID,
			User:     *user,
			Title:    title,
			Name:     name,
			Public:   public,
			Subtitle: subtitle,
		},
	})
}

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

func postCourseReleasesFORM(c *gin.Context) {
	courseID := c.Params.ByName("courseID")

	course, err := db.GetCourseWithIDPreloadUser(courseID)
	if err != nil {
		log.Println("v2/course.go ERROR getting course:", err)
		c.JSON(http.StatusInternalServerError, payload{
			Error: "Error getting course",
		})
		return
	}

	releases, _ := db.GetAnyReleases(course.ID)

	num := uint16(1)
	if len(releases) > 0 {
		num = uint16(len(releases) + 1)
	}

	// if author then they can view private releases
	release := db.Release{
		Num:      num,
		CourseID: course.ID,
	}
	err1 := db.CreateRelease(&release)
	if err1 != nil {
		log.Println("v2/course.go ERROR creating release in postCourseReleasesFORM:", err1)
		c.JSON(http.StatusInternalServerError, payload{
			Error: "Error creating release",
		})
		return
	}

	c.JSON(http.StatusOK, payload{
		Payload: release,
	})
}
