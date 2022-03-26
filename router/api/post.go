package api

import (
	"log"
	"main/db"
	"main/router/session"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func courseVersionNewPost(c *gin.Context) {
	if !session.IsLoggedInValid(c) {
		log.Println("api LOG not logged in valid")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := session.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api ERROR couldn't get logged in user:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	versionID := c.Params.ByName("versionID")
	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api ERROR getting version:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	/*
		free courses do not need to be "purchased"
		until you want to actually start writing a post for it
		below follows the logic for this
		- if user has purchased the course
			- if course price is zero
				- give free course purchase!
			- else prevent from posting, since they can't have access
		- else prevent from posting, since they can't have access
	*/

	// check if user has purchased the course
	if !db.UserHasPurchasedCourse(user.ID, version.ReleaseID) {
		release, err4 := db.GetRelease(version.ReleaseID)
		if err4 != nil {
			log.Println("api ERROR getting release:", err4)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// if course price is zero, give free course purchase!
		if release.Price == 0 {
			purchase := db.Purchase{
				UserID:     user.ID,
				ReleaseID:  release.ID,
				AmountPaid: 0,
			}

			err5 := db.CreatePurchase(&purchase)
			if err5 != nil {
				log.Println("api ERROR creating purchase:", err5)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("api LOG course is not free")
			// else prevent from posting, since they can't have access
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	sectionID := c.PostForm("sectionID")
	markdown := c.PostForm("markdown")

	// prevent from posting empty posts
	if strings.Trim(markdown, " ") == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	post := db.Post{
		UserID:   user.ID,
		Markdown: markdown,
	}
	err2 := db.CreatePost(&post)
	if err2 != nil {
		log.Println("api ERROR creating post:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sectionIDNum, err4 := strconv.ParseUint(sectionID, 10, 64)
	if err4 != nil {
		log.Println("api ERROR parsing sectionID:", err4)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	postToRelease := db.PostToRelease{
		PostID:    post.ID,
		ReleaseID: version.ReleaseID,
		SectionID: sectionIDNum,
	}
	err3 := db.CreatePostToRelease(&postToRelease)
	if err3 != nil {
		log.Println("api ERROR creating postToRelease:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println("posted new post:", post)
	log.Println("markdown is:", markdown)
}

func postUpdatePost(c *gin.Context) {
	if !session.IsLoggedInValid(c) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := session.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api ERROR couldn't get logged in user:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	postID := c.Params.ByName("postID")

	post, err := db.GetPost(postID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if post.UserID != user.ID {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newMarkdown := c.PostForm("markdown")

	if strings.Trim(newMarkdown, " ") == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db.UpdatePost(postID, newMarkdown)
}

func postEditSectionContent(c *gin.Context) {
	log.Println("editing content")
}
