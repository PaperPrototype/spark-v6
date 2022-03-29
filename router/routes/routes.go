// backend html page routes
package routes

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	// landing page
	router.GET("/", getLanding) // index

	// user page
	router.GET("/:username", getUser)         // get user profile
	router.GET("/:username/media/:mediaName") // where the user can upload and access images or gifs
	/*
		user setting on a cog wheel button, but don't offer settings menu as a url route
	*/

	// course landing pages
	router.GET("/:username/:course", getCourse)                    // course page, gives newest release
	router.GET("/:username/:course/:releaseNum", getCourseRelease) // course page, gives specific release

	// privately accessed course settings
	router.GET("/:username/:course/settings", mustBeCourseEditor, getCourseSettings)
	router.POST("/:username/:course/settings/display", mustBeCourseEditor, postCourseSettingsDisplay)
	router.POST("/:username/:course/settings/release/new", mustBeCourseEditor, postNewRelease)
	router.GET("/:username/:course/settings/release/delete", mustBeCourseEditor, getReleaseDelete)
	router.POST("/:username/:course/settings/release/delete/confirm", mustBeCourseEditor, postReleaseDeleteConfirm)
	router.POST("/:username/:course/settings/release/edit", mustBeCourseEditor, postEditRelease)
	router.POST("/:username/:course/settings/version/new", mustBeCourseEditor, postNewVersion)
	router.POST("/:username/:course/settings/version/delete", mustBeCourseEditor, postDeleteVersion)

	// view inside of course content
	router.GET("/:username/:course/view/:versionID", MustHaveAccessToCourseRelease, getCourseVersion)                   // view a version of the course
	router.GET("/:username/:course/view/:versionID/:sectionID", MustHaveAccessToCourseRelease, getCourseVersionSection) // view a section of the course
	router.GET("/:username/:course/view/:versionID/posts")                                                              // view posts
	router.GET("/:username/:course/view/:versionID/posts/:postID")                                                      // view specific post
	router.GET("/:username/:course/view/:versionID/posts/user/:username")                                               // view posts by a specific user
	router.GET("/:username/:course/view/:versionID/chat")                                                               // view the live chatroom

	// payments routes for courses
	router.POST("/:username/:course/buy/:releaseID", mustBeLoggedIn, postBuyRelease)
	router.GET("/:username/:course/success", mustBeLoggedIn, getBuySuccess)
	router.GET("/:username/:course/cancel", mustBeLoggedIn, getBuyCancel)

	// course media assets (zip, png, gif)
	router.GET("/media/:versionID/name/:mediaName", getNameMedia)
	router.GET("/media/:versionID/id/:mediaID")

	// auth
	router.GET("/signup", getSignup) // make a new account
	router.POST("/signup", postSignup)
	router.GET("/login", getLogin)   // log into existing account
	router.POST("/login", postLogin) // signinto account
	router.GET("/login/verify")      // verify account
	router.GET("/login/forgot")      // forgot password
	router.GET("/logout", getLogout) // logout

	// Payouts (for logged in users only)
	router.GET("/user/payouts", mustBeLoggedIn, getUserPayouts)            // get user payouts
	router.GET("/user/payouts/connect", mustBeLoggedIn, getPayoutsConnect) // connect account to stripe so we can pay out to teachers
	router.GET("/user/payouts/refresh", mustBeLoggedIn, getPayoutsRefresh)
	router.GET("/user/payouts/connect/return", mustBeLoggedIn, getPayoutsConnectFinished)
	router.GET("/user/payouts/payout", mustBeLoggedIn, getPayout)

	router.GET("/courses", getCourses) // search courses with possible url query
	/*
		/courses?query=Intro+to+coding&order=relevance
	*/

	// creating a new course
	router.GET("/new", getNew)
	router.POST("/new", postNew)

	router.GET("/lost", getLost)
}
