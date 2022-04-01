package routes

import (
	"log"
	"main/db"
	"main/helpers"
	"main/msg"
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
)

func getSettings(c *gin.Context) {
	user := auth.GetLoggedInUserLogError(c)

	c.HTML(
		http.StatusOK,
		"settings.html",
		gin.H{
			"Menu":     "General",
			"User":     user,
			"Messages": msg.GetMessages(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getSettingsCourses(c *gin.Context) {
	user := auth.GetLoggedInUserLogError(c)

	courses, err := user.GetAuthorCourses()
	if err != nil {
		log.Println("routes/settings ERROR getting AuthorCourses in getSettingsCourses:", err)
	}

	c.HTML(
		http.StatusOK,
		"settingsCourses.html",
		gin.H{
			"Courses":  courses,
			"Menu":     "Courses",
			"User":     user,
			"Messages": msg.GetMessages(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getStripeConnect(c *gin.Context) {
	user, err2 := auth.GetLoggedInUser(c)
	if err2 != nil {
		log.Println("routes/payments ERROR getting logged in user:", err2)
		msg.SendMessage(c, "Error getting user.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	if !user.Verified {
		msg.SendMessage(c, "Error getting user.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	if user.HasStripeConnection() {
		// user already has account connection
		// they just need to "link" aka input info for their account with stripe
		c.Redirect(http.StatusFound, "/settings/stripe/connect/refresh")
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
		c.Redirect(http.StatusFound, "/settings/payouts")
		return
	}

	// create stripeConnection in db with expressAccount.ID
	stripeConnection := db.StripeConnection{
		StripeAccountID: expressAccount.ID,
		UserID:          user.ID,
	}
	err3 := db.CreateStripeConnection(&stripeConnection)
	if err3 != nil {
		log.Println("routes/payments ERROR creating stripeConnection in db:", err3)
		msg.SendMessage(c, "Error connecting accouint to stripe.")
		c.Redirect(http.StatusFound, "/settings/courses")
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
		RefreshURL: stripe.String(helpers.GetHost() + "/settings/stripe/connect/refresh"),

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
		ReturnURL: stripe.String(helpers.GetHost() + "/settings/stripe/connect/return"),
		Type:      stripe.String("account_onboarding"),
	}

	accountLink, err1 := accountlink.New(params2)
	if err1 != nil {
		log.Println("routes/payments ERROR creating connected account link:", err1)
		msg.SendMessage(c, "Error creating connected account link.")
		c.Redirect(http.StatusFound, "/settings/payouts")
		return
	}

	c.Redirect(http.StatusSeeOther, accountLink.URL)
}

func getStripeRefresh(c *gin.Context) {
	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes/payments ERROR getting user for getPayoutsRefresh:", err)
		msg.SendMessage(c, "Error getting user")
		c.Redirect(http.StatusFound, "/settings/payouts")
		return
	}

	stripeConnection, err2 := db.GetStripeConnection(user.ID)
	if err2 != nil {
		log.Println("routes/payments ERROR getting stripe connection for getPayoutsRefresh:", err2)
		msg.SendMessage(c, "Error getting stripe connection")
		c.Redirect(http.StatusFound, "/settings/payouts")
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
		RefreshURL: stripe.String(helpers.GetHost() + "/settings/stripe/connect/refresh"),

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
		ReturnURL: stripe.String(helpers.GetHost() + "/settings/stripe/connect/return"),
		Type:      stripe.String("account_onboarding"),
	}

	accountLink, err1 := accountlink.New(params2)
	if err1 != nil {
		log.Println("routes/payments ERROR creating connected account link getPayoutsRefresh:", err1)
		msg.SendMessage(c, "Error creating connected account link.")
		c.Redirect(http.StatusFound, "/settings/courses")
		return
	}

	c.Redirect(http.StatusSeeOther, accountLink.URL)
}

func getStripeConnectFinished(c *gin.Context) {
	// TODO test if user successfully connected in stripe (check the state of the details_submitted parameter)
	// see https://stripe.com/docs/connect/express-accounts#return_url
	// code example https://stripe.com/docs/api/accounts/retrieve

	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes/payments ERROR getting user in getPayoutsConnectFinished:", err)
		msg.SendMessage(c, "You must be logged in!")
		c.Redirect(http.StatusFound, "/")
		return
	}

	stripeConnection, err1 := db.GetStripeConnection(user.ID)
	if err1 != nil {
		log.Println("routes/payments ERROR getting stripeConnection in getPayoutsConnectFinished:", err1)
		msg.SendMessage(c, "Error getting stripe connection")
		c.Redirect(http.StatusFound, "/settings/courses")
		return
	}

	// if all details were submitted
	submitted, err2 := stripeConnection.DetailsSubmitted()
	if err2 != nil {
		log.Println("routes/payments ERROR getting details submitted:", err2)
		msg.SendMessage(c, "Error getting account details.")
	}

	if !submitted {
		msg.SendMessage(c, "Finish filling out account details by clicking 'Connect account' again. Make sure to use the same email.")
		c.Redirect(http.StatusFound, "/settings/courses")
		return
	}

	// user is successfully connected
	msg.SendMessage(c, "Successfully connected account!")
	// user is successfully connected
	c.Redirect(http.StatusFound, "/settings/courses")
}
