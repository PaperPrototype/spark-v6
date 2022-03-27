// backend html page routes
package routes

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	// landing page
	router.GET("/", getLanding) // index

	// course landing pages
	router.GET("/:course", getCourse)                    // course page, gives newest release
	router.GET("/:course/:releaseNum", getCourseRelease) // course page, gives specific release

	// course settings
	router.GET("/:course/settings", mustBeCourseEditor, getCourseSettings)
	router.POST("/:course/settings/display", mustBeCourseEditor, postCourseSettingsDisplay)
	router.POST("/:course/settings/release/new", mustBeCourseEditor, postNewRelease)
	router.GET("/:course/settings/release/delete", mustBeCourseEditor, getReleaseDelete)
	router.POST("/:course/settings/release/delete/confirm", mustBeCourseEditor, postReleaseDeleteConfirm)
	router.POST("/:course/settings/release/edit", mustBeCourseEditor, postEditRelease)
	router.POST("/:course/settings/version/new", mustBeCourseEditor, postNewVersion)
	router.POST("/:course/settings/version/delete", mustBeCourseEditor, postDeleteVersion)

	// view inside of course content
	router.GET("/:course/view/:versionID", MustHaveAccessToCourseRelease, getCourseVersion)                   // view a version of the course
	router.GET("/:course/view/:versionID/:sectionID", MustHaveAccessToCourseRelease, getCourseVersionSection) // view a section of the course
	router.GET("/:course/view/:versionID/posts")                                                              // view posts
	router.GET("/:course/view/:versionID/posts/:postID")                                                      // view specific post
	router.GET("/:course/view/:versionID/posts/user/:username")                                               // view posts by a specific user
	router.GET("/:course/view/:versionID/chat")                                                               // view the live chatroom

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

	// publicly accessible pages for logged in users
	router.GET("/user/:username", getUser)         // get user profile
	router.GET("/user/:username/media/:mediaName") // where the user can upload and access images or gifs
	/*
		user setting on a cog wheel button, but don't offer settings menu as a url route
	*/

	// DASHBOARD for logged in users only
	router.GET("/user/courses", mustBeLoggedIn)                 // get user payouts
	router.GET("/user/payouts", mustBeLoggedIn, getUserPayouts) // get user payouts
	router.GET("/user/payouts/connect", mustBeLoggedIn)         // connect account to stripe so we can pay out to teachers

	router.GET("/courses", getCourses) // search courses with possible url query
	/*
		/courses?query=Intro+to+coding&order=relevance
	*/

	// creating a new course
	router.GET("/new", getNew)
	router.POST("/new", postNew)

	router.GET("/lost", getLost)

	// payments routes for courses
	router.POST("/:course/buy/:releaseID", mustBeLoggedIn, postBuyRelease)
	router.GET("/:course/success", mustBeLoggedIn, getBuySuccess)
	router.GET("/:course/cancel", mustBeLoggedIn, getBuyCancel)
}
