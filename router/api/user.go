package api

import (
	"log"
	"main/db"
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserPosts(c *gin.Context) {
	username := c.Params.ByName("username")

	user, err := db.GetUserWithUsername(username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("api/user.go ERROR getting user in getUserPosts:", err)
		return
	}

	posts, err1 := db.GetUserPosts(user.ID)
	if err1 != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("api/user.go ERROR getting posts in getUserPosts:", err1)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func getUserCourses(c *gin.Context) {
	username := c.Params.ByName("username")

	user, err := db.GetUserWithUsername(username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("api/user.go ERROR getting user in getUserCourses:", err)
		return
	}

	courses, err1 := user.GetPublicPurchasedCourses()
	if err1 != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("api/user.go ERROR getting public purchased courses in getUserCourses:", err1)
		return
	}

	c.JSON(http.StatusOK, courses)
}

// if the user is logged in and looking at their own profile then also show private courses
func getUserAuthoredCourses(c *gin.Context) {
	username := c.Params.ByName("username")

	user, err := db.GetUserWithUsername(username)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("api/user.go ERROR getting user in getUserAuthoredCourses:", err)
		return
	}

	if auth.IsLoggedInValid(c) {
		loggedInUser := auth.GetLoggedInUserLogError(c)
		if user.ID == loggedInUser.ID {
			courses, err1 := user.GetPublicAndPrivateAuthoredCourses()
			if err1 != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				log.Println("api/user.go ERROR getting PublicAndPrivateAuthoredCourses in getUserAuthoredCourses:", err1)
				return
			}

			c.JSON(http.StatusOK, courses)
			return
		}
	}

	courses, err1 := user.GetPublicAuthoredCourses()
	if err1 != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println("api/user.go ERROR getting PublicAuthoredCourses in getUserAuthoredCourses:", err1)
		return
	}

	c.JSON(http.StatusInternalServerError, courses)
}
