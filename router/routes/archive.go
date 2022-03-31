package routes

import (
	"log"
	"main/db"
	"main/msg"
	"main/router/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

// archived route function

func getUserPayouts(c *gin.Context) {
	if !session.IsLoggedInValid(c) {
		msg.SendMessage(c, "You must be logged in to access payouts.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	user, err := session.GetLoggedInUser(c)
	if err != nil {
		msg.SendMessage(c, "Error getting logged in user")
		c.Redirect(http.StatusFound, "/")
		return
	}

	courses, err1 := db.GetUserCourses(user.ID)
	if err1 != nil {
		log.Println("routes ERROR getting user courses from getUserPayouts:", err1)
	}

	var totalPayout float64 = 0
	for _, course := range courses {
		totalPayout += course.GetCurrentTotalCoursePayoutAmountLogError()
	}

	stripeConnection, err2 := db.GetStripeConnection(user.ID)
	if err2 != nil {
		log.Println("routes/get ERROR getting stripe connection for getUserPayouts:", err2)
	}

	c.HTML(
		http.StatusOK,
		"payout.html",
		gin.H{
			"TotalPayout":      totalPayout,
			"StripeConnection": stripeConnection,
			"Courses":          courses,
			"Messages":         msg.GetMessages(c),
			"User":             user,
			"LoggedIn":         session.IsLoggedInValid(c),
			"Meta":             metaDefault,
		},
	)
}
