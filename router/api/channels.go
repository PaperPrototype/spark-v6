package api

import (
	"log"
	"main/db"
	"main/helpers"
	"main/router/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ChannelPayload struct {
	Messages []db.Message
	Channel  db.Channel
}

// send a message
func postChannelSendMessage(c *gin.Context) {
	channelID := c.Params.ByName("channelID")
	markdown := c.PostForm("markdown")
	versionID := c.Params.ByName("versionID")

	if markdown == "" || markdown == " " {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	version, err3 := db.GetVersion(versionID)
	if err3 != nil {
		log.Println("api/channels ERROR getting version in postChannelSendMessage:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	course, err4 := db.GetCoursePreloadUser(version.CourseID)
	if err4 != nil {
		log.Println("api/channels ERROR getting course in postChannelSendMessage:", err4)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	channelIDNum, err := strconv.ParseUint(channelID, 10, 64)
	if err != nil {
		log.Println("api/channels ERROR parsing channelID in postChannelNewMessage:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := auth.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api/channels ERROR getting logged in user in postChannelNewMessage:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	message := db.Message{
		UserID:    user.ID,
		ChannelID: channelIDNum,
		Markdown:  markdown,
	}
	err2 := db.CreateMessage(&message)
	if err2 != nil {
		log.Println("api/channels ERROR creating message in postChannelNewMessage:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	usernames := helpers.GetUserMentions(markdown)
	err5 := db.NotifyUsers(usernames, "@"+user.Username+" mentioned you in the chat of "+course.Title, "/"+course.User.Username+"/"+course.Name+"/view/"+versionID+"?channel_id="+channelID)
	if err5 != nil {
		log.Println("api/channels ERROR notifying users in postChannelNewMessage:", err5)
	}
}

func getChannelMessages(c *gin.Context) {
	channelID := c.Params.ByName("channelID")
	newest := c.Query("newest")

	channel, err3 := db.GetChannel(channelID)
	if err3 != nil {
		log.Println("db/channels ERROR getting channel in getChannelMessages:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// start timeout handler
	start := time.Now()

	// comments requested!

	// if the user does not know which is the newest comment
	// loading initial comments
	if newest == "" {
		// get the initial comments
		// limit to last 20 comments
		messages, count, err := db.GetMessages(channelID, 20)
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

			// get the initial messages
			// limit to last 20 messsages
			messages, count, err = db.GetMessages(channelID, 20)
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
			ChannelPayload{
				Messages: messages,
				Channel:  *channel,
			},
		)
		return
	}

	// otherwise the user knows the newest message date
	messages, count, err1 := db.GetNewMessages(channelID, newest)
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
		messages, count, err1 = db.GetNewMessages(channelID, newest)
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
		ChannelPayload{
			Messages: messages,
			Channel:  *channel,
		},
	)
}
