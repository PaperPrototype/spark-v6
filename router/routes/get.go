package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"main/conn"
	"main/db"
	"main/mailer"
	"main/markdown"
	"main/msg"
	"main/router/auth"

	"github.com/gin-gonic/gin"
)

var metaDefault = Meta{
	Title: "Sparker, Epic Software Courses",
	Desc:  "Time to ditch software degree's and switch to portfolio's",
}

func getCourse(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	course, err := db.GetUserCoursePreloadUser(username, courseName)
	if err != nil {
		log.Println("ERROR getting course:", err)
		notFound(c)
		return
	}

	release, err1 := db.GetNewestPublicCourseRelease(course.ID)
	if err1 != nil {
		log.Println("ERROR getting release:", err1)

		// render without release
		c.HTML(
			http.StatusOK,
			"course.html",
			gin.H{
				"Course":   course,
				"Messages": msg.GetMessages(c),
				"User":     auth.GetLoggedInUserLogError(c),
				"LoggedIn": auth.IsLoggedInValid(c),
				"Meta":     metaDefault,
			},
		)
		return
	}

	// user has clicked on course from searching, skip the landing page if it is free and take them straight to the course
	if release.Price == 0 {
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/view/"+fmt.Sprint(release.GetNewestVersionLogError().ID))
		return
	}

	// convert release desc to support markdown
	releaseMarkdowned, err5 := markdown.Convert([]byte(release.Markdown))
	if err5 != nil {
		log.Println("routes/get course ERROR converting markown for release Desc:", err5)
	}

	release.Markdown = template.HTML(releaseMarkdowned.String())

	purchased := false

	if release.Price != 0 {
		user, err2 := auth.GetLoggedInUser(c)
		if err2 == nil {
			if user.HasPurchasedRelease(release.ID) {
				purchased = true
			}
		}
	}

	c.HTML(
		http.StatusOK,
		"course.html",
		gin.H{
			"Purchased": purchased,
			"Course":    course,
			"Release":   release,
			"Messages":  msg.GetMessages(c),
			"User":      auth.GetLoggedInUserLogError(c),
			"LoggedIn":  auth.IsLoggedInValid(c),
			"Meta":      metaDefault,
		},
	)
}

func getCourses(c *gin.Context) {
	search, _ := c.GetQuery("query")
	sort, _ := c.GetQuery("sort")

	c.HTML(
		http.StatusOK,
		"courses.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Search":   search,
			"Sort":     sort,
			"Meta":     metaDefault,
		},
	)
}

func getLanding(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"landing.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getLogin(c *gin.Context) {
	redirectURL := c.Query("redirect_url")

	if redirectURL != "" {
		log.Println("redirect url is:", redirectURL)
	}

	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"RedirectURL": redirectURL,
			"Messages":    msg.GetMessages(c),
			"User":        auth.GetLoggedInUserLogError(c),
			"LoggedIn":    auth.IsLoggedInValid(c),
			"Meta":        metaDefault,
		},
	)
}

func getNew(c *gin.Context) {
	if !auth.IsLoggedInValid(c) {
		msg.SendMessage(c, "You must be logged in to create a new course.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	user, err2 := auth.GetLoggedInUser(c)
	if err2 != nil {
		log.Println("ERROR getting logged in user:", err2)
		msg.SendMessage(c, "Error getting logged in user.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	if !user.Verified {
		msg.SendMessage(c, "You must verify your account before you can upload courses.")
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	if !user.HasStripeConnection() {
		msg.SendMessage(c, "You must connect your account to stripe before you can upload courses.")
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	c.HTML(
		http.StatusOK,
		"new.html",
		gin.H{
			"Course":   db.Course{},
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getSignup(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getLogout(c *gin.Context) {
	auth.Logout(c)
	c.Redirect(http.StatusFound, "/")
}

func getCourseSettings(c *gin.Context) {
	if !auth.IsLoggedInValid(c) {
		msg.SendMessage(c, "You must be logged in to access a settings page.")
		notFound(c)
		return
	}

	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	course, err := db.GetUserCoursePreloadUser(username, courseName)
	if err != nil {
		log.Println("ERROR getting course:", err)
		notFound(c)
		return
	}

	user, err1 := db.GetUser(course.UserID)
	if err1 != nil {
		log.Println("ERROR getting user:", err1)
		notFound(c)
		return
	}

	c.HTML(
		http.StatusOK,
		"courseSettings.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     user,
			"LoggedIn": auth.IsLoggedInValid(c),
			"Course":   course,
			"Releases": course.GetAllCourseReleasesLogError(),
			"Meta":     metaDefault,
		},
	)
}

func getLost(c *gin.Context) {
	c.HTML(
		http.StatusNotFound,
		"notFound.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getCourseVersion(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	// TODO check that user owns course release before allowing access

	course, err1 := db.GetUserCoursePreloadUser(username, courseName)
	if err1 != nil {
		log.Println("routes ERROR getting course from db:", err1)
		notFound(c)
		return
	}

	version, err := course.GetVersion(versionID)
	if err != nil {
		log.Println("routes ERROR getting version from db:", err)
		msg.SendMessage(c, "That course upload may have been deleted.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	release, err3 := db.GetAllRelease(version.ReleaseID)
	if err3 != nil {
		log.Println("routes ERROR getting release from db:", err3)
		msg.SendMessage(c, "That course release is not available.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	section := &db.Section{}
	var err2 error
	if !release.HasGithubRelease() {
		section, err2 = version.GetFirstSectionPreload()
		if err2 != nil {
			log.Println("routes/get ERROR getting versions first section in getCourseVersion:", err2)
		}
	}

	var progress int64

	if auth.IsLoggedInValid(c) {
		user := auth.GetLoggedInUserLogError(c)

		amount := course.GetNewestPublicCourseReleaseLogError().UserPostsCountLogError(user.ID)
		total := version.SectionsCountLogError()

		log.Println("amount:", amount)
		log.Println("total:", total)

		// convert to float for division
		floatProgress := float64(amount) / float64(total)

		// convert deciaml to percentage
		floatProgress *= 100

		// cast
		progress = int64(floatProgress)

		log.Println("progress:", progress)
	}

	c.HTML(
		http.StatusOK,
		"courseView.html",
		gin.H{
			"Release":  release,
			"Course":   course,
			"Version":  version,
			"Section":  section,
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
			"Progress": progress,
		},
	)
}

func getCourseRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	releaseNum := c.Params.ByName("releaseNum")

	course, err := db.GetUserCoursePreloadUser(username, courseName)
	if err != nil {
		log.Println("ERROR getting course:", err)
		notFound(c)
		return
	}

	release, err1 := db.GetCourseReleaseNumString(course.ID, releaseNum)
	if err1 != nil {
		log.Println("ERROR getting release:", err1)

		// render without release
		c.HTML(
			http.StatusOK,
			"course.html",
			gin.H{
				"Course":   course,
				"Messages": msg.GetMessages(c),
				"User":     auth.GetLoggedInUserLogError(c),
				"LoggedIn": auth.IsLoggedInValid(c),
				"Meta":     metaDefault,
			},
		)
		return
	}

	// convert release desc to support markdown
	releaseMarkdowned, err5 := markdown.Convert([]byte(release.Markdown))
	if err5 != nil {
		log.Println("routes/get course ERROR converting markown for release Desc:", err5)
	}

	release.Markdown = template.HTML(releaseMarkdowned.String())

	purchased := false

	if release.Price != 0 {
		user, err2 := auth.GetLoggedInUser(c)
		if err2 == nil {
			if user.HasPurchasedRelease(release.ID) {
				purchased = true
			}
		}
	}

	c.HTML(
		http.StatusOK,
		"course.html",
		gin.H{
			"Purchased": purchased,
			"Course":    course,
			"Release":   release,
			"Messages":  msg.GetMessages(c),
			"User":      auth.GetLoggedInUserLogError(c),
			"LoggedIn":  auth.IsLoggedInValid(c),
			"Meta":      metaDefault,
		},
	)
}

func getCourseVersionSection(c *gin.Context) {
	courseName := c.Params.ByName("course")
	username := c.Params.ByName("username")
	versionID := c.Params.ByName("versionID")

	// TODO check that user owns course release before allowing access

	course, err1 := db.GetUserCoursePreloadUser(username, courseName)
	if err1 != nil {
		log.Println("routes ERROR getting course from db:", err1)
		notFound(c)
		return
	}

	version, err := course.GetVersion(versionID)
	if err != nil {
		log.Println("routes ERROR getting version from db:", err)
		notFound(c)
		return
	}

	release, err3 := db.GetAllRelease(version.ReleaseID)
	if err3 != nil {
		log.Println("routes ERROR getting release from db:", err3)
		msg.SendMessage(c, "That course release is not available.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName)
		return
	}

	sectionID := c.Params.ByName("sectionID")
	section := &db.Section{}
	var err2 error
	if !release.HasGithubRelease() {
		section, err2 = db.GetSectionPreload(sectionID)
		if err2 != nil {
			log.Println("routes/get ERROR getting versions first section in getCourseVersion:", err2)
			msg.SendMessage(c, "Error getting first section.")
		}
	}

	var progress int64

	if auth.IsLoggedInValid(c) {
		user := auth.GetLoggedInUserLogError(c)

		amount := course.GetNewestPublicCourseReleaseLogError().UserPostsCountLogError(user.ID)
		total := version.SectionsCountLogError()

		log.Println("amount:", amount)
		log.Println("total:", total)

		// convert to float for division
		floatProgress := float64(amount) / float64(total)

		// convert deciaml to percentage
		floatProgress *= 100

		// cast
		progress = int64(floatProgress)

		log.Println("progress:", progress)
	}

	c.HTML(
		http.StatusOK,
		"courseView.html",
		gin.H{
			"Release":  release,
			"Course":   course,
			"Version":  version,
			"Section":  section,
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
			"Progress": progress,
		},
	)
}

func getUser(c *gin.Context) {
	username := c.Params.ByName("username")

	if username == "" {
		msg.SendMessage(c, "Can't find that user!")
		notFound(c)
		return
	}

	profileUser, err := db.GetUserWithUsername(username)
	if err != nil {
		log.Println("routes ERROR gettingUserWithUsername:", err)
		msg.SendMessage(c, "Can't find that user!")
		notFound(c)
		return
	}

	courses, err1 := profileUser.GetPurchasedCourses()
	if err1 != nil {
		log.Println("routes ERROR getting courses for user:", err1)
	}

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"Messages":       msg.GetMessages(c),
			"User":           auth.GetLoggedInUserLogError(c),
			"ProfileUser":    profileUser,
			"ProfileCourses": courses,
			"LoggedIn":       auth.IsLoggedInValid(c),
			"Meta":           metaDefault,
		},
	)
}

func getNameMedia(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	mediaName := c.Params.ByName("mediaName")
	media, err := db.GetMedia(versionID, mediaName)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Writer.Header().Set("Content-Type", media.Type)
	c.Writer.Header().Set("Content-Length", fmt.Sprint(media.Length))

	conn := conn.GetConn()
	WriteMediaChunks(conn, c.Writer, media.ID)
}

func getReleaseDelete(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	releaseID := c.Query("releaseID")

	release, err := db.GetAllRelease(releaseID)
	if err != nil {
		log.Println("routes/getReleaseDelete ERROR getting release:", err)
	}

	c.HTML(
		http.StatusOK,
		"confirmDelete.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,

			// special params for confirmDelete.html
			"Action":  "/" + username + "/" + courseName + "/settings/release/delete/confirm",
			"Message": "Confirm you want to delete release " + fmt.Sprint(release.Num),
			"Data":    release.ID,
			"Further": "This will also delete all versions and user content in this release!",
		},
	)
}

func getVerify(c *gin.Context) {
	verifyUUID := c.Params.ByName("verifyUUID")
	verify, err := db.GetVerify(verifyUUID)
	if err != nil {
		log.Println("routes/get ERROR getting verify in getVerify:", err)
		msg.SendMessage(c, "Error or link has expired.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	user, err1 := db.GetUser(verify.UserID)
	if err1 != nil {
		log.Println("routes/get ERROR getting user in getVerify:", err)
		msg.SendMessage(c, "Error user.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	message := "Failed to verify user"
	err2 := user.SetVerified(true)
	if err2 != nil {
		log.Println("routes/get ERROR setting verified to true in getVerify:", err2)
	} else {
		message = "You've been verified!"
	}

	c.HTML(
		http.StatusOK,
		"verify.html",
		gin.H{
			"Message":  message,
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getNewVerify(c *gin.Context) {
	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes/get ERROR getting logged in user in getNewVerify:", err)
		msg.SendMessage(c, "Error getting logged in user")
		c.Redirect(http.StatusFound, "/")
		return
	}

	err1 := mailer.SendVerification(user.ID)
	if err1 != nil {
		log.Println("routes/get ERROR sending verification email in getNewVerify:", err1)
		msg.SendMessage(c, "Error sending verification email")
		c.Redirect(http.StatusFound, "/")
		return
	}

	msg.SendMessage(c, "Sent verification link to your email. Make sure to check your spam folder.")
	c.Redirect(http.StatusFound, "/settings")
}

func getAbout(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"about.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}
