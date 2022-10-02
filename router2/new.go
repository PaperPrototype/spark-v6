package router2

import (
	"log"
	"main/auth2"
	"main/db"
	"main/helpers"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"course_new.html",
		gin.H{
			"User":     auth2.GetLoggedInUserLogError(c),
			"Messages": msg.GetMessages(c),
			"Meta": meta{
				Title: "Sparker - New Course",
			},
		},
	)
}

func postNew(c *gin.Context) {
	name := c.PostForm("name")
	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")

	if title == "" {
		msg.SendMessage(c, "The tile must not be empty")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	if name == "" {
		msg.SendMessage(c, "The URL name must not be empty")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	if !helpers.IsAllowedUsername(name) {
		msg.SendMessage(c, "The URL name can only contain the following lowercase characters: abcdefghijklmnopqrstuvwxyz1234567890-_")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	user := auth2.GetLoggedInUserLogError(c)

	if !user.Verified {
		msg.SendMessage(c, "You must verify your email before you can upload courses.")
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	available, err1 := db.UserCourseNameAvailable(user.Username, name)
	if err1 != nil {
		log.Println("router2/new.go ERROR checking if course url name is available in postNew:", err1)
		msg.SendMessage(c, "That url name is already taken.")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	if !available {
		msg.SendMessage(c, "That url name is already taken.")
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	course := db.Course{
		Name:     name,
		Title:    title,
		Subtitle: subtitle,
		UserID:   user.ID,
	}
	err := db.CreateCourse(&course)
	if err != nil {
		msg.SendMessage(c, "Error creating course")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	msg.SendMessage(c, "Welcome to your new course! You can start by making a new release!")
	c.Redirect(http.StatusFound, "/"+user.Username+"/"+name)
}
