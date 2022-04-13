package api

import (
	"log"
	"main/db"
	"main/router/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NotifsPayload struct {
	Notifs []db.Notif
	Count  int64
}

func getNewNotifications(c *gin.Context) {
	newest := c.Query("newest")

	user := auth.GetLoggedInUserLogError(c)

	// start timing for request timeout
	start := time.Now()

	// loading initial notifications
	if newest == "" {
		// get the initial notifications
		// limit to last 20 notifications
		notifs, count, err := db.GetUnreadNotifs(user.ID, 5)
		if err != nil {
			log.Println("api/get ERROR getting initial notifications in getNotifications:", err)
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

			// get the initial notifications
			// limit to last 20 notiifications
			notifs, count, err = db.GetUnreadNotifs(user.ID, 5)
			if err != nil {
				log.Println("api/get ERROR getting initial notifications in getNotifications:", err)
				c.Status(http.StatusInternalServerError)
				return
			}

			// if there are new notifications
			if count > 0 {
				break
			}
		}

		// sending initial notifications!
		c.JSON(
			http.StatusOK,
			NotifsPayload{
				Notifs: notifs,
				Count:  count,
			},
		)
		return
	}

	// otherwise the user knows the newest notifications
	notifs, count, err1 := db.GetNewUnreadNotifs(user.ID, newest)
	if err1 != nil {
		// this cause the comment system to send another query
		log.Println("api/get ERROR getting new notifications in getNotifications:", err1)
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

		// check for new notifications
		notifs, count, err1 = db.GetNewUnreadNotifs(user.ID, newest)
		if err1 != nil {
			// this will cause the comments system to send another query
			log.Println("api/get ERROR getting new notifications in getNotifications:", err1)
			c.Status(http.StatusInternalServerError)
			return
		}

		// if there are new notifications
		if count > 0 {
			break
		}
	}

	c.JSON(
		http.StatusOK,
		NotifsPayload{
			Notifs: notifs,
			Count:  count,
		},
	)
}

func postDoneNotification(c *gin.Context) {
	err := db.NotifSetRead(c.PostForm("notifID"))
	if err != nil {
		log.Println("api/notifs ERROR setting notification as read:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// tell js frontend everything is ok and worked
	c.Writer.WriteHeader(http.StatusOK)
}
