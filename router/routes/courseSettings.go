package routes

import (
	"log"
	"main/db"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postSettingsNewPrerequisite(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	preqCourseID := c.PostForm("preqCourseID")

	course, err := db.GetUserCoursePreload(username, courseName)
	if err != nil {
		log.Println("api/courseSettings ERROR getting course in postSettingsPrerequisitesAdd:", err)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"settings")
		return
	}

	preqCourse, err1 := db.GetCoursePreloadUser(preqCourseID)
	if err1 != nil {
		log.Println("api/courseSettings ERROR getting prerequisite course in postSettingsPrerequisitesAdd:", err1)
		msg.SendMessage(c, "Error getting prerequisite course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	preq := db.Prerequisite{
		CourseID:             course.ID,
		PrerequisiteCourseID: preqCourse.ID,
	}

	err2 := db.CreatePrerequisite(&preq)
	if err2 != nil {
		log.Println("api/courseSettings ERROR creating prerequisite in postSettingsPrerequisitesAdd:", err2)
		msg.SendMessage(c, "Error creating prerequisite.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	err3 := db.UpdateCourseLevel(course.ID)
	if err3 != nil {
		log.Println("api/courseSettings ERROR updating course level in postSettingsPrerequisitesAdd:", err3)
		msg.SendMessage(c, "An error occurred while updating the course level")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	err4 := db.UpdateDependantCourseLevels(course.ID)
	if err4 != nil {
		log.Println("api/courseSettings ERROR updating dependant course levels in postSettingsPrerequisitesAdd:", err4)
		msg.SendMessage(c, "An error occurred while updating the courses that depend on this course")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully added prerequisite course.")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postSettingsRemovePrerequisite(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	preqID := c.PostForm("preqID")

	err := db.DeletePrerequisite(preqID)
	if err != nil {
		log.Println("routes/courseSettings ERROR dleting prerequisite in postSettingsRemovePrerequisite:", err)
		msg.SendMessage(c, "Error deleting prerequisite course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	course, err := db.GetUserCoursePreload(username, courseName)
	if err != nil {
		log.Println("api/courseSettings ERROR getting course in postSettingsRemovePrerequisite:", err)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"settings")
		return
	}

	err3 := db.UpdateCourseLevel(course.ID)
	if err3 != nil {
		log.Println("api/courseSettings ERROR updating course level in postSettingsRemovePrerequisite:", err3)
		msg.SendMessage(c, "An error occurred while updating the course level")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	err4 := db.UpdateDependantCourseLevels(course.ID)
	if err4 != nil {
		log.Println("api/courseSettings ERROR updating dependant course levels in postSettingsRemovePrerequisite:", err4)
		msg.SendMessage(c, "An error occurred while updating the courses that depend on this course")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully deleted prerequisite.")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}
