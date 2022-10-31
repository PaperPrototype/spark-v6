package router2

import (
	"main/middlewares"
)

func SetupRoutes() {
	// pahge to make new course
	router.GET("/new", middlewares.MustBeLoggedIn, getNew)
	router.POST("/new", middlewares.MustBeLoggedIn, postNew)

	// landing and homepage
	router.GET("/", getBrowse)

	// 404 page
	router.NoRoute(getLost)

	// get course
	router.GET("/:username/:course", getCourse)
	router.GET("/:username/:course/:sectionID", getCourse)

	// TODO use course.ID instead of username and courseName
	// incase name changes when user clicks buy
	router.GET("/:username/:course/buy/:releaseID", getBuyRelease)
	router.GET("/:username/:course/buy/success", middlewares.MustBeLoggedIn, getBuySuccess)
	router.GET("/:username/:course/buy/cancel", middlewares.MustBeLoggedIn, getBuyCancel)

	// get github based media

	// auth
	router.POST("/login", postLogin)
	router.POST("/signup", postSignup)
	router.GET("/logout", getLogout)

	router.GET("/login/verify/:verifyUUID", getVerify) // verify account
	router.GET("/login/verify/new", getNewVerify)      // send verification email

	// settings
	router.GET("/settings", getSettings)
	router.POST("/settings/edit/user", middlewares.MustBeLoggedIn, postSettingsEditUser)
	router.POST("/settings/edit/email", middlewares.MustBeLoggedIn, postSettingsEditEmail)

	router.GET("/settings/github/connect", middlewares.MustBeLoggedIn, getGithubConnect)
	router.GET("/settings/github/connect/return", middlewares.MustBeLoggedIn, getGithubConnectFinished)

	router.GET("/settings/stripe/connect", middlewares.MustBeLoggedIn, getStripeConnect)
	router.GET("/settings/stripe/login", middlewares.MustBeLoggedIn, getStripeLogin)
	router.GET("/settings/stripe/connect/refresh", middlewares.MustBeLoggedIn, getStripeRefresh)
	router.GET("/settings/stripe/connect/return", middlewares.MustBeLoggedIn, getStripeConnectFinished)

	// stripe payments
	router.POST("/stripe/webhooks", postStripeWebhook)
}
