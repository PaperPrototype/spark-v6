package router2

import (
	"fmt"
	"log"
	"main/auth2"
	"main/db"
	"main/helpers"
	"main/mailer"
	"main/msg"
	"main/payments"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/checkout/session"
)

// 1 day and 5 minutes?
const PaymentExpiresAfter time.Duration = (24 * time.Hour) + (time.Minute * 5)

func getBuyRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	releaseID := c.Params.ByName("releaseID")

	if !auth2.IsLoggedInValid(c) {
		msg.SendMessage(c, "You must be logged in to purchase a course")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	release, err3 := db.GetPublicReleaseWithID(releaseID)
	if err3 != nil {
		log.Println("routes/payments ERROR getting release:", err3)
		msg.SendMessage(c, "Error getting course release")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	user, err5 := auth2.GetLoggedInUser(c)
	if err5 != nil {
		msg.SendMessage(c, "You must be logged in to access this page.")

		return
	}

	if user.OwnsRelease(release.ID) {
		msg.SendMessage(c, "You already own this course release!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	course, err := db.GetUserCoursePreload(username, courseName)
	if err != nil {
		log.Println("routes/payments ERROR getting course:", err)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	author, err6 := db.GetUser(course.UserID)
	if err6 != nil {
		log.Println("routes/payments ERROR getting course author:", err6)
		msg.SendMessage(c, "Error getting course author.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	if !author.HasStripeConnection() {
		msg.SendMessage(c, "An error has occured! The author of this course cannot accept payments at this time. We apologize. We are gifting you the course for free :)")

		mailer.SendStripePaymentProblemEmail(author.ID, "It appears someone tried to purchase a course from you, but you don't have stripe setup (the online payments service we use). We gifted the course for free since you could not accept the payment.")

		ownership := db.Ownership{
			Desc:      payments.DescStripeConnectionNotSetup,
			UserID:    user.ID,
			CourseID:  release.CourseID,
			ReleaseID: release.ID,
			Completed: false,
		}
		err5 := db.CreateOwnership(&ownership)
		if err5 != nil {
			msg.SendMessage(c, "Failed to create course ownership! That is not supposed to happen! Contact us and send a screenshot of this message!")
			c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
			return
		}

		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	if course.UserID == user.ID {
		msg.SendMessage(c, "You are the author of this course!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	stripeConnection, err8 := db.GetStripeConnection(author.ID)
	if err8 != nil {
		log.Println("toures/payments ERROR getting stripe connection postBuyRelease:", err8)
		msg.SendMessage(c, "There was an error. But we'll gift you the course for free :)")

		mailer.SendStripePaymentProblemEmail(author.ID, "There was an error getting your stripe connection! Your stripe info may need updated! We had to gift the course for free since you could not accept the payment.")

		ownership := db.Ownership{
			Desc:      payments.DescStripeConnectionNotSetupError,
			UserID:    user.ID,
			CourseID:  release.CourseID,
			ReleaseID: release.ID,
			Completed: false,
		}
		err5 := db.CreateOwnership(&ownership)
		if err5 != nil {
			msg.SendMessage(c, "Failed to create course ownership! That is not supposed to happen! Contact us and send a screenshot of this message!")
			c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
			return
		}

		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	payoutsEnabled, err9 := stripeConnection.PayoutsEnabled()
	if err9 != nil {
		log.Println("toures/payments ERROR getting stripe connection postBuyRelease:", err8)
		msg.SendMessage(c, "There was an error. But we'll gift you the course for free :)")

		mailer.SendStripePaymentProblemEmail(author.ID, "There was an error checking if yu can accept payouts! Your stripe info may need updated! We gifted the course for free since you could not accept the payment.")

		ownership := db.Ownership{
			Desc:      payments.DescStripeChargesNotEnabledError,
			UserID:    user.ID,
			CourseID:  release.CourseID,
			ReleaseID: release.ID,
			Completed: false,
		}
		err5 := db.CreateOwnership(&ownership)
		if err5 != nil {
			msg.SendMessage(c, "Failed to create course ownership! That is not supposed to happen! Contact us and send a screenshot of this message!")
			c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
			return
		}

		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	// allow payment attempt to go through, but send a warning message!
	if !payoutsEnabled {
		mailer.SendStripePaymentProblemEmail(author.ID, "Your stripe info needs updated. We had to gift one of your courses for free since your stripe account could not accept payments at this time.")
	}

	// release.Price * PercentageShare
	// uint16        * float32
	var sparkersCut int64 = int64(math.Round(float64(float32(release.Price) * payments.PercentageShare)))

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(string(stripe.CurrencyUSD)),
					UnitAmount: stripe.Int64(int64(release.Price)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(course.Title + " v" + fmt.Sprint(release.Num)),
						Description: stripe.String(course.Subtitle),
					},
				},
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(sparkersCut),
			Description:          stripe.String("Buying " + course.Title + " v" + fmt.Sprint(release.Num)),
			TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
				Destination: stripe.String(stripeConnection.StripeAccountID),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(helpers.GetHost() + "/" + username + "/" + courseName + "/buy/success"),
		CancelURL:  stripe.String(helpers.GetHost() + "/" + username + "/" + courseName + "/buy/cancel"),
	}

	resultSession, err7 := session.New(params)
	if err7 != nil {
		log.Println("routes/payments ERROR creating checkout session in postBuyRelease:", err7)
		msg.SendMessage(c, "Error creating payment session.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	buyRelease := db.AttemptBuyRelease{
		StripeSessionID: resultSession.ID,
		ReleaseID:       release.ID,
		UserID:          user.ID,
		AmountPaying:    release.Price,
		ExpiresAt:       time.Now().Add(PaymentExpiresAfter),
	}
	err4 := db.CreateBuyRelease(&buyRelease)
	if err4 != nil {
		log.Println("routes/payments ERROR creating buyRelease:", err4)
		msg.SendMessage(c, "Error creating buyRelease")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	c.Redirect(http.StatusSeeOther, resultSession.URL)
}
