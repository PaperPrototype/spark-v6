package routes

import (
	"fmt"
	"log"
	"main/db"
	"main/helpers"
	"main/mailer"
	"main/msg"
	"main/payments"
	"math"
	"net/http"
	"time"

	"main/router/auth"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

/*
resources

start here:
	- https://stripe.com/docs/payments/accept-a-payment#create-product-prices-upfront
then here for confirming payments:
	- https://stripe.com/docs/payments/checkout/custom-success-page#modify-success-url

then for payouts modify code a little to have a TransferData (aka transfer funds to ceonnected account):
	- https://stripe.com/docs/connect/collect-then-transfer-guide
*/

const PaymentExpiresAfter time.Duration = (24 * time.Hour) + (time.Minute * 5)

func warnUserStripeConnectionProblem(author *db.User, messageProblem string) {

	err10 := db.NotifyUsers([]string{author.Username}, `<i class="fa-solid fa-triangle-exclamation"></i> IMPORTANT: `+messageProblem, "/settings/teaching")
	if err10 != nil {
		log.Println("routes/payments ERROR sending notification to user @"+author.Username+" in postBuyRelease:", err10)
	}

	mailer.SendStripePaymentProblemEmail(author.ID, messageProblem)
}

func postBuyRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	releaseID := c.Params.ByName("releaseID")
	courseName := c.Params.ByName("course")

	release, err3 := db.GetPublicReleaseWithID(releaseID)
	if err3 != nil {
		log.Println("routes/payments ERROR getting release:", err3)
		msg.SendMessage(c, "Error getting course release")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	user, err5 := auth.GetLoggedInUser(c)
	if err5 != nil {
		msg.SendMessage(c, "You must be logged in to access this page.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	if user.HasPurchasedRelease(release.ID) {
		msg.SendMessage(c, "You already own this course release!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
		return
	}

	course, err := db.GetUserCoursePreloadUser(username, courseName)
	if err != nil {
		log.Println("routes/payments ERROR getting course:", err)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	author, err6 := db.GetUser(course.UserID)
	if err6 != nil {
		log.Println("routes/payments ERROR getting course author:", err6)
		msg.SendMessage(c, "Error getting course author.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	if !author.HasStripeConnection() {
		msg.SendMessage(c, "The author of this course cannot accept payments at this time. We'll gift you the course for free :)")

		warnUserStripeConnectionProblem(author, "It appears someone tried to purchase a course from you, but you don't have a stripe connection! We had to gift the course for free since you could not accept payments at this time.")

		purchase := db.Purchase{
			VersionID:  release.GetNewestVersionLogError().ID,
			UserID:     user.ID,
			ReleaseID:  release.ID,
			CourseID:   release.CourseID,
			Desc:       payments.DescStripeConnectionNotSetup,
			AmountPaid: 0,
			AuthorsCut: 0,
			CreatedAt:  time.Now(),
		}

		err7 := db.CreatePurchase(&purchase)
		if err7 != nil {
			log.Println("routes/payments ERROR creating purchase:", err7)
			msg.SendMessage(c, "Error gifting course.")
		}

		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	if course.UserID == user.ID {
		msg.SendMessage(c, "You are the author of this course!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
		return
	}

	stripeConnection, err8 := db.GetStripeConnection(author.ID)
	if err8 != nil {
		log.Println("toures/payments ERROR getting stripe connection postBuyRelease:", err8)
		msg.SendMessage(c, "There was an error. But we'll gift you the course for free :)")

		warnUserStripeConnectionProblem(author, "There was an error getting your stripe connection! Your stripe info may need updated! We had to gift the course for free since you could not accept the payment.")

		purchase := db.Purchase{
			VersionID:  release.GetNewestVersionLogError().ID,
			UserID:     user.ID,
			ReleaseID:  release.ID,
			CourseID:   release.CourseID,
			Desc:       payments.DescStripeConnectionNotSetupError,
			AmountPaid: 0,
			AuthorsCut: 0,
			CreatedAt:  time.Now(),
		}
		err7 := db.CreatePurchase(&purchase)
		if err7 != nil {
			log.Println("routes/payments ERROR creating purchase:", err7)
			msg.SendMessage(c, "Error gifting course.")
		}

		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	chargesEnabled, err9 := stripeConnection.ChargesEnabled()
	if err9 != nil {
		log.Println("toures/payments ERROR getting stripe connection postBuyRelease:", err8)
		msg.SendMessage(c, "There was an error. But we'll gift you the course for free :)")

		warnUserStripeConnectionProblem(author, "There was an error getting your stripe connection! Your stripe info may need updated! We had to gift the course for free since you could not accept the payment.")

		purchase := db.Purchase{
			VersionID:  release.GetNewestVersionLogError().ID,
			UserID:     user.ID,
			ReleaseID:  release.ID,
			CourseID:   release.CourseID,
			Desc:       payments.DescStripeChargesNotEnabledError,
			AmountPaid: 0,
			AuthorsCut: 0,
			CreatedAt:  time.Now(),
		}
		err7 := db.CreatePurchase(&purchase)
		if err7 != nil {
			log.Println("routes/payments ERROR creating purchase:", err7)
			msg.SendMessage(c, "Error gifting course.")
		}

		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	// allow payment attempt to go through, but send a warning message!
	if !chargesEnabled {
		// msg.SendMessage(c, "The author of this course cannot accept payments at this time. We'll gift you the course for free :)")

		warnUserStripeConnectionProblem(author, "Your stripe info needs updated. We had to gift one of your courses for free since your stripe account could not accept payments at this time.")

		// purchase := db.Purchase{
		// 	VersionID:  release.GetNewestVersionLogError().ID,
		// 	UserID:     user.ID,
		// 	ReleaseID:  release.ID,
		// 	CourseID:   release.CourseID,
		// 	Desc:       payments.DescStripeChargesNotEnabled,
		// 	AmountPaid: 0,
		// 	AuthorsCut: 0,
		// 	CreatedAt:  time.Now(),
		// }

		// err7 := db.CreatePurchase(&purchase)
		// if err7 != nil {
		// 	log.Println("routes/payments ERROR creating purchase:", err7)
		// 	msg.SendMessage(c, "Error gifting course.")
		// }

		// c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		// return
	}

	// release.Price * PercentageShare
	// uint16        * float32
	var sparkersCut int64 = int64(math.Round(float64(float32(release.Price) * payments.PercentageShare)))

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Name:        stripe.String(course.Title),
				Amount:      stripe.Int64(int64(release.Price)),
				Currency:    stripe.String(string(stripe.CurrencyUSD)),
				Description: stripe.String("Buying " + course.Name + " v" + fmt.Sprint(release.Num)),
				Quantity:    stripe.Int64(1),
			},
		},
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ApplicationFeeAmount: stripe.Int64(sparkersCut),
			Description:          stripe.String("Buying " + course.Name + " v" + fmt.Sprint(release.Num)),
			TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
				Destination: stripe.String(stripeConnection.StripeAccountID),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(helpers.GetHost() + "/" + username + "/" + courseName + "/buy/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(helpers.GetHost() + "/" + username + "/" + courseName + "/buy/cancel"),
	}

	resultSession, err7 := session.New(params)
	if err7 != nil {
		log.Println("routes/payments ERROR creating checkout session in postBuyRelease:", err7)
		msg.SendMessage(c, "Error creating payment session.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	buyRelease := db.AttemptBuyRelease{
		StripeSessionID: resultSession.ID,
		StripePaymentID: resultSession.PaymentIntent.ID,
		ReleaseID:       release.ID,
		UserID:          user.ID,
		AmountPaying:    release.Price,
		ExpiresAt:       time.Now().Add(PaymentExpiresAfter),
	}
	err4 := db.CreateBuyRelease(&buyRelease)
	if err4 != nil {
		log.Println("routes/payments ERROR creating buyRelease:", err4)
		msg.SendMessage(c, "Error creating buyRelease")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	c.Redirect(http.StatusSeeOther, resultSession.URL)
}

func getBuySuccess(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	// get special id and check if we have saved it in the db and is a valid payment
	stripeSessionID := c.Query("session_id")

	if stripeSessionID == "" {
		log.Println("db MALICIOUS behaviour trying to success a order that was never started or expired?")
		msg.SendMessage(c, "An error occurred.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	buyRelease, err := db.GetBuyRelease(stripeSessionID)
	if err != nil {
		log.Println("db ERROR getting buyRelease:", err)
		log.Println("db MALICIOUS behaviour trying to success a order that was never started or expired?")
		msg.SendMessage(c, "An error occurred. Maybe your order timed out?")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
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
		// if user gets signed out while buying make sure they can log in again and still finish getting their purchased course!
		c.BindQuery(&struct {
			RedirectURL string `form:"redirect_url"`
		}{
			RedirectURL: "/" + username + "/" + courseName + "buy?session_id=" + stripeSessionID,
		})
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if user.ID != buyRelease.UserID {
		msg.SendMessage(c, "This is not your purchase!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	// AmountPayed * 0.15
	// we took 15 percent
	ourCut := uint64(math.Round(float64(float32(buyRelease.AmountPaying) * payments.PercentageShare)))
	authorsCut := buyRelease.AmountPaying - ourCut

	purchase := db.Purchase{
		UserID:          user.ID,
		VersionID:       version.ID,
		ReleaseID:       buyRelease.ReleaseID,
		StripeSessionID: buyRelease.StripeSessionID,
		StripePaymentID: buyRelease.StripePaymentID,
		CourseID:        version.CourseID,
		CreatedAt:       time.Now(),
		AmountPaid:      buyRelease.AmountPaying,
		AuthorsCut:      authorsCut,
	}
	err2 := db.CreatePurchase(&purchase)
	if err2 != nil {
		msg.SendMessage(c, "Purchase creating failed! That is not supposed to happen! Contact us and send a screenshot of this message!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	// if purchase creation succeeded
	// delete buyRelease to prevent user from re-buying for free
	err4 := db.DeleteBuyRelease(buyRelease.StripeSessionID)
	if err4 != nil {
		log.Println("failed to delete buyRelease (may have timed out and auto deleted):", err4)
	}

	log.Println("Payment success!")
	msg.SendMessage(c, "Payment success! Welcome!")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/view/"+fmt.Sprint(version.ID))
}

func getBuyCancel(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	// delete buyRelease if user cancels payment
	stripeSessionID := c.Query("session_id")
	if stripeSessionID != "" {
		err := db.DeleteBuyRelease(stripeSessionID)
		if err != nil {
			log.Println("routes/payments.go ERROR failed to delete buyRelease (may have timed out and auto deleted):", err)
		}
	}

	log.Println("Payment canceled!")
	msg.SendMessage(c, "Payment canceled")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
}
