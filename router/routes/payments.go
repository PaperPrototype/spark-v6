package routes

import (
	"fmt"
	"log"
	"main/db"
	"main/helpers"
	"main/msg"
	"net/http"
	"time"

	auth "main/router/session"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

func postBuyRelease(c *gin.Context) {
	stripe.Key = helpers.GetStripeKey()

	releaseID := c.Params.ByName("releaseID")
	courseName := c.Params.ByName("course")

	user, err5 := auth.GetLoggedInUser(c)
	if err5 != nil {
		msg.SendMessage(c, "You must be logged in to access this page.")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	course, err := db.GetCourse(courseName)
	if err != nil {
		log.Println("routes/payments ERROR getting course:", err)
		msg.SendMessage(c, "Error getting course")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	release, err3 := db.GetPublicReleaseWithIDStr(releaseID)
	if err3 != nil {
		log.Println("routes/payments ERROR getting release:", err3)
		msg.SendMessage(c, "Error getting course release")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	if course.UserID == user.ID {
		msg.SendMessage(c, "You are the author of this course!")
		c.Redirect(http.StatusFound, "/"+course.Name+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
		return
	}

	if user.HasPurchasedRelease(release.ID) {
		msg.SendMessage(c, "You already own this course release!")
		c.Redirect(http.StatusFound, "/"+course.Name+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
		return
	}

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(course.Title),
					},
					UnitAmount: stripe.Int64(int64(release.Price) * 100),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:8080/" + courseName + "/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:8080/" + courseName + "/cancel"),
	}

	resultSession, err2 := session.New(params)
	if err2 != nil {
		log.Println("routes/payments ERROR creating payment session:", err2)
		msg.SendMessage(c, "Error creating payment session")
		c.Redirect(http.StatusFound, "/"+course.Name)
		return
	}

	buyRelease := db.BuyRelease{
		ID:           resultSession.ID,
		ReleaseID:    release.ID,
		UserID:       user.ID,
		AmountPaying: release.Price,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	err4 := db.CreateBuyRelease(&buyRelease)
	if err4 != nil {
		log.Println("routes/payments ERROR creating buyRelease:", err4)
		msg.SendMessage(c, "Error creating buyRelease")
		c.Redirect(http.StatusFound, "/"+course.Name)
		return
	}

	c.Redirect(http.StatusSeeOther, resultSession.URL)
}

func getBuySuccess(c *gin.Context) {
	courseName := c.Params.ByName("course")

	// get special id and check if we have saved it in the db and is a valid payment
	buyReleaseID := c.Query("session_id")

	if buyReleaseID == "" {
		log.Println("db MALICIOUS behavviour trying to success a order that was never started or expired?")
		msg.SendMessage(c, "An error occurred.")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	buyRelease, err := db.GetBuyRelease(buyReleaseID)
	if err != nil {
		log.Println("db ERROR getting buyRelease:", err)
		log.Println("db MALICIOUS behaviour trying to success a order that was never started or expired?")
		msg.SendMessage(c, "An error occurred. Maybe your order timed out?")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	version, err3 := db.GetNewestReleaseVersion(buyRelease.ReleaseID)
	if err3 != nil {
		// this error may happen if user chooses to buy a course with no versions released yet
		log.Println("Error getting version:", err3)
		msg.SendMessage(c, "No course versions released yet. Don't worry once this author releases a version you will have access to it!")
		// continue DON'T RETURN!
	}

	user, err1 := auth.GetLoggedInUser(c)
	if err1 != nil {
		msg.SendMessage(c, "An error occurred. Are you logged in?")
		// if user gets signed out while buy make sure they can log in again and still finish purchasing course!
		c.BindQuery(&struct {
			RedirectURL string `form:"redirect_url"`
		}{
			RedirectURL: "/" + courseName + "buy?session_id=" + buyReleaseID,
		})
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if user.ID != buyRelease.UserID {
		msg.SendMessage(c, "This is not your purchase!")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	purchase := db.Purchase{
		UserID:     user.ID,
		VersionID:  version.ID,
		ReleaseID:  buyRelease.ReleaseID,
		CreatedAt:  time.Now(),
		AmountPaid: buyRelease.AmountPaying,
	}
	err2 := db.CreatePurchase(&purchase)
	if err2 != nil {
		msg.SendMessage(c, "Purchase creating failed! That is not supposed to happen! Contact us and send a screenshot of this message!")
		c.Redirect(http.StatusFound, "")
		return
	}

	// if purchase creation succeeded
	// delete buyRelease to prevent user from re-buying for free
	err4 := db.DeleteBuyRelease(buyRelease.ID)
	if err4 != nil {
		log.Println("failed to delete buyRelease (may have timed out and auto deleted):", err4)
	}

	log.Println("Payment success!")
	msg.SendMessage(c, "Payment success! Welcome!")
	c.Redirect(http.StatusFound, "/"+courseName+"/view/"+fmt.Sprint(version.ID))
}

func getBuyCancel(c *gin.Context) {
	courseName := c.Params.ByName("course")

	log.Println("Payment canceled!")
	msg.SendMessage(c, "Payment canceled")
	c.Redirect(http.StatusFound, "/"+courseName)
}
