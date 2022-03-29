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
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	"github.com/stripe/stripe-go/v72/checkout/session"
)

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

	course, err := db.GetUserCoursePreloadUser(username, courseName)
	if err != nil {
		log.Println("routes/payments ERROR getting course:", err)
		msg.SendMessage(c, "Error getting course")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	user, err5 := auth.GetLoggedInUser(c)
	if err5 != nil {
		msg.SendMessage(c, "You must be logged in to access this page.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	if course.UserID == user.ID {
		msg.SendMessage(c, "You are the author of this course!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
		return
	}

	if user.HasPurchasedRelease(release.ID) {
		msg.SendMessage(c, "You already own this course release!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
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
		SuccessURL: stripe.String(helpers.GetHost() + "/" + username + "/" + courseName + "/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(helpers.GetHost() + "/" + username + "/" + courseName + "/cancel"),
	}

	resultSession, err2 := session.New(params)
	if err2 != nil {
		log.Println("routes/payments ERROR creating payment session:", err2)
		msg.SendMessage(c, "Error creating payment session")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/"+fmt.Sprint(release.Num))
		return
	}

	buyRelease := db.BuyRelease{
		StripeSessionID: resultSession.ID,
		ReleaseID:       release.ID,
		UserID:          user.ID,
		AmountPaying:    release.Price,
		ExpiresAt:       time.Now().Add(24 * time.Hour),
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
	buyReleaseID := c.Query("session_id")

	if buyReleaseID == "" {
		log.Println("db MALICIOUS behavviour trying to success a order that was never started or expired?")
		msg.SendMessage(c, "An error occurred.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	buyRelease, err := db.GetBuyRelease(buyReleaseID)
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
			RedirectURL: "/" + username + "/" + courseName + "buy?session_id=" + buyReleaseID,
		})
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if user.ID != buyRelease.UserID {
		msg.SendMessage(c, "This is not your purchase!")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	purchase := db.Purchase{
		UserID:     user.ID,
		VersionID:  version.ID,
		ReleaseID:  buyRelease.ReleaseID,
		CourseID:   version.CourseID,
		CreatedAt:  time.Now(),
		AmountPaid: buyRelease.AmountPaying,
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

	log.Println("Payment canceled!")
	msg.SendMessage(c, "Payment canceled")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
}

func getPayoutsConnect(c *gin.Context) {
	user, err2 := auth.GetLoggedInUser(c)
	if err2 != nil {
		log.Println("routes/payments ERROR getting logged in user:", err2)
		msg.SendMessage(c, "Error getting user.")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	if user.HasStripeConnection() {
		// user already has account connection
		// they just need to "link" aka input info for their account with stripe
		c.Redirect(http.StatusFound, "/user/payouts/refresh")
		return
	}

	// see full params list and examples here https://stripe.com/docs/connect/express-accounts
	params := &stripe.AccountParams{
		Type: stripe.String(string(stripe.AccountTypeExpress)),
	}

	expressAccount, err := account.New(params)
	if err != nil {
		log.Println("routes/payments ERROR creating connected stripe account:", err)
		msg.SendMessage(c, "Error creating connected account.")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	// TODO
	// create stripeConnection in db with expressAccount.ID
	stripeConnection := db.StripeConnection{
		StripeAccountID: expressAccount.ID,
		UserID:          user.ID,
	}
	err3 := db.CreateStripeConnection(&stripeConnection)
	if err3 != nil {
		log.Println("routes/payments ERROR creating stripeConnection in db:", err3)
		msg.SendMessage(c, "Error connecting accouint to stripe.")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	params2 := &stripe.AccountLinkParams{
		Account: stripe.String(expressAccount.ID),

		/*
			Your user is redirected to the refresh_url when:

			- The link has expired (a few minutes have passed since the link was created).
			- The link was already visited (the user refreshed the page or clicked back or forward in their browser).
			- Your platform is no longer able to access the account.
			- The account has been rejected.
			The refresh_url should call Account Links again on your server with the same parameters and redirect the user to the Connect Onboarding flow to create a seamless experience.
		*/
		RefreshURL: stripe.String(helpers.GetHost() + "/user/payouts/refresh"),

		/*
			Stripe issues a redirect to this URL when the user completes the Connect
			Onboarding flow. This doesn’t mean that all information has been collected
			or that there are no outstanding requirements on the account. This only
			means the flow was entered and exited properly.

			No state is passed through this URL. After a user is redirected to your
			return_url, check the state of the details_submitted parameter on their
			account by doing either of the following:

			- Listening to account.updated events.
			- Calling the Accounts API (with expressAccount.ID) and inspecting the returned object.
		*/
		ReturnURL: stripe.String(helpers.GetHost() + "/user/payouts/connect/return"),
		Type:      stripe.String("account_onboarding"),
	}

	accountLink, err1 := accountlink.New(params2)
	if err1 != nil {
		log.Println("routes/payments ERROR creating connected account link:", err1)
		msg.SendMessage(c, "Error creating connected account link.")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	c.Redirect(http.StatusSeeOther, accountLink.URL)
}

func getPayoutsRefresh(c *gin.Context) {
	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes/payments ERROR getting user for getPayoutsRefresh:", err)
		msg.SendMessage(c, "Error getting user")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	stripeConnection, err2 := db.GetStripeConnection(user.ID)
	if err2 != nil {
		log.Println("routes/payments ERROR getting stripe connection for getPayoutsRefresh:", err2)
		msg.SendMessage(c, "Error getting stripe connection")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	params2 := &stripe.AccountLinkParams{
		Account: stripe.String(stripeConnection.StripeAccountID),

		/*
			Your user is redirected to the refresh_url when:

			- The link has expired (a few minutes have passed since the link was created).
			- The link was already visited (the user refreshed the page or clicked back or forward in their browser).
			- Your platform is no longer able to access the account.
			- The account has been rejected.
			The refresh_url should call Account Links again on your server with the same parameters and redirect the user to the Connect Onboarding flow to create a seamless experience.
		*/
		RefreshURL: stripe.String(helpers.GetHost() + "/user/payouts/refresh"),

		/*
			Stripe issues a redirect to this URL when the user completes the Connect
			Onboarding flow. This doesn’t mean that all information has been collected
			or that there are no outstanding requirements on the account. This only
			means the flow was entered and exited properly.

			No state is passed through this URL. After a user is redirected to your
			return_url, check the state of the details_submitted parameter on their
			account by doing either of the following:

			- Listening to account.updated events.
			- Calling the Accounts API (with expressAccount.ID) and inspecting the returned object.
		*/
		ReturnURL: stripe.String(helpers.GetHost() + "/user/payouts/connect/return"),
		Type:      stripe.String("account_onboarding"),
	}

	accountLink, err1 := accountlink.New(params2)
	if err1 != nil {
		log.Println("routes/payments ERROR creating connected account link getPayoutsRefresh:", err1)
		msg.SendMessage(c, "Error creating connected account link.")
		c.Redirect(http.StatusFound, "/user/payouts")
		return
	}

	c.Redirect(http.StatusSeeOther, accountLink.URL)
}

func getPayoutsConnectFinished(c *gin.Context) {
	// TODO test if user successfully connected in stripe (check the state of the details_submitted parameter)
	// see https://stripe.com/docs/connect/express-accounts#return_url
	// code example https://stripe.com/docs/api/accounts/retrieve
	log.Println("/user/payouts/connect/return")

	// user is successfully connected
	msg.SendMessage(c, "Successfully connected account!")
	// user is successfully connected
	c.Redirect(http.StatusFound, "/user/payouts")
}
