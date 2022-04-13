package api

import (
	"fmt"
	"log"
	"main/db"
	"main/helpers"
	"main/router/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func postPostComment(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	postID := c.Params.ByName("postID")
	markdown := c.PostForm("markdown")

	postIDNum, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		log.Println("api/post ERROR parsing uint in postPostComment:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := auth.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api/post ERROR getting user in postPostComment:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create the comment
	comment := db.Comment{
		UserID:   user.ID,
		PostID:   postIDNum,
		Markdown: markdown,
	}
	err2 := db.CreateComment(&comment)
	if err2 != nil {
		log.Println("api/post ERROR creating comment in postPostComment:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	version, err3 := db.GetVersion(versionID)
	if err3 != nil {
		log.Println("api/comments ERROR getting version in postPostComment:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	author, err4 := version.GetAuthorUser()
	if err4 != nil {
		log.Println("api/comments ERROR getting version author user in postPostComment:", err4)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	course, err5 := db.GetCourse(version.CourseID)
	if err5 != nil {
		log.Println("api/comments ERROR getting version course in postPostComment:", err5)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	usernames := helpers.GetUserMentions(markdown)
	log.Println("usernames to notify:", usernames)
	err6 := db.NotifyUsers(usernames, "@"+user.Username+" mentioned you in a post's comments in "+course.Title, "/"+author.Username+"/"+course.Name+"/view/"+fmt.Sprint(version.ID)+"?post_id="+postID)
	if err6 != nil {
		log.Println("api/comments ERROR notifying users in postPostComment:", err6)
	}

	c.Status(http.StatusOK)
}

func getPostComments(c *gin.Context) {
	postID := c.Params.ByName("postID")
	newest := c.Query("newest")

	// start timeout handler
	start := time.Now()

	// comments requested!

	// if the user does not know which is the newest comment
	// loading initial comments
	if newest == "" {
		// get the initial comments
		// limit to last 20 comments
		comments, count, err := db.GetComments(postID, 20)
		if err != nil {
			log.Println("api/get ERROR getting initial comments in getPostComments:", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		for count == 0 {
			// if nothing found then sleep 3 seconds before querying again
			time.Sleep(SleepTime * time.Second)

			// have timout so if user disconnects we don't keep going forever.
			// Otherwise we will query even after user disconnects!!!
			duration := time.Since(start)

			// only try for up to 20 seconds
			if duration.Seconds() > MaxTimeoutSeconds {
				break
			}

			// get the initial comments
			// limit to last 20 comments
			comments, count, err = db.GetComments(postID, 20)
			if err != nil {
				log.Println("api/get ERROR getting initial comments in getPostComments:", err)
				c.Status(http.StatusInternalServerError)
				return
			}

			// if there are new comments
			if count > 0 {
				break
			}
		}

		// sending initial comments!
		c.JSON(
			http.StatusOK,
			comments,
		)
		return
	}

	// otherwise the user knows the newest comment
	comments, count, err1 := db.GetNewComments(postID, newest)
	if err1 != nil {
		// this cause the comment system to send another query
		log.Println("api/get ERROR getting new comments in getPostComments:", err1)
		c.Status(http.StatusInternalServerError)
		return
	}

	// while count == 0
	for count == 0 {
		// if nothing found then sleep 3 seconds before querying again
		time.Sleep(SleepTime * time.Second)

		// have timout so if user disconnects we don't keep going forever.
		// Otherwise we will query even after user disconnects!!!
		duration := time.Since(start)

		// only try for up to 20 seconds
		if duration.Seconds() > MaxTimeoutSeconds {
			break
		}

		// check for new comments
		comments, count, err1 = db.GetNewComments(postID, newest)
		if err1 != nil {
			// this will cause the comments system to send another query
			log.Println("api/get ERROR getting new comments in getPostComments:", err1)
			c.Status(http.StatusInternalServerError)
			return
		}

		// if there are new comments
		if count > 0 {
			break
		}
	}

	// sending new comments!
	c.JSON(
		http.StatusOK,
		comments,
	)
}
