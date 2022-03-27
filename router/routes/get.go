package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"main/conn"
	"main/db"
	"main/markdown"
	"main/msg"
	"main/router/session"

	"github.com/gin-gonic/gin"
)

var metaDefault = Meta{
	Title: "Spark, Epic software courses",
	Desc:  "Time to ditch software degree's and switch to portfolio's",
}

func getCourse(c *gin.Context) {
	name := c.Params.ByName("course")

	course, err := db.GetCourse(name)
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
				"User":     session.GetLoggedInUserHideError(c),
				"LoggedIn": session.IsLoggedInValid(c),
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
		user, err2 := session.GetLoggedInUser(c)
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
			"User":      session.GetLoggedInUserHideError(c),
			"LoggedIn":  session.IsLoggedInValid(c),
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
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
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
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
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
			"User":        session.GetLoggedInUserHideError(c),
			"LoggedIn":    session.IsLoggedInValid(c),
			"Meta":        metaDefault,
		},
	)
}

func getNew(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"new.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
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
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getLogout(c *gin.Context) {
	session.Logout(c)
	c.Redirect(http.StatusFound, "/")
}

func getCourseSettings(c *gin.Context) {
	if !session.IsLoggedInValid(c) {
		msg.SendMessage(c, "You must be logged in to access a settings page.")
		notFound(c)
		return
	}

	course, err := db.GetCourse(c.Params.ByName("course"))
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
			"LoggedIn": session.IsLoggedInValid(c),
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
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getCourseVersion(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	courseName := c.Params.ByName("course")

	// TODO check that user owns course release before allowing access

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("routes ERROR getting version from db:", err)
		msg.SendMessage(c, "No course version yet!")
		c.Redirect(http.StatusFound, "/"+courseName)
		return
	}

	section, err2 := version.GetFirstSectionPreload()
	if err2 != nil {
		msg.SendMessage(c, "Error getting first section.")
	}

	course, err1 := db.GetCoursewithID(version.CourseID)
	if err1 != nil {
		log.Println("routes ERROR getting course from db:", err1)
		notFound(c)
		return
	}

	var progress int64

	if session.IsLoggedInValid(c) {
		user, _ := session.GetLoggedInUser(c)

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
			"Course":   course,
			"Version":  version,
			"Section":  section,
			"Messages": msg.GetMessages(c),
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
			"Meta":     metaDefault,
			"Progress": progress,
		},
	)
}

func getCourseRelease(c *gin.Context) {
	name := c.Params.ByName("course")
	releaseNum := c.Params.ByName("releaseNum")

	course, err := db.GetCourse(name)
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
				"User":     session.GetLoggedInUserHideError(c),
				"LoggedIn": session.IsLoggedInValid(c),
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
		user, err2 := session.GetLoggedInUser(c)
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
			"User":      session.GetLoggedInUserHideError(c),
			"LoggedIn":  session.IsLoggedInValid(c),
			"Meta":      metaDefault,
		},
	)
}

func getCourseVersionSection(c *gin.Context) {
	versionID := c.Params.ByName("versionID")

	// TODO check that user owns course release before allowing access

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("routes ERROR getting version from db:", err)
		notFound(c)
		return
	}

	sectionID := c.Params.ByName("sectionID")

	section, err2 := db.GetSectionPreload(sectionID)
	if err2 != nil {
		log.Println("routes/get ERROR getting section for getCourseVersionSection:", err2)
		msg.SendMessage(c, "Error getting first section.")
	}

	course, err1 := db.GetCoursewithID(version.CourseID)
	if err1 != nil {
		log.Println("routes ERROR getting course from db:", err1)
		notFound(c)
		return
	}

	var progress int64

	if session.IsLoggedInValid(c) {
		user, _ := session.GetLoggedInUser(c)

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
			"Course":   course,
			"Version":  version,
			"Section":  section,
			"Messages": msg.GetMessages(c),
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
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

	courses, err1 := profileUser.GetCourses()
	if err1 != nil {
		log.Println("routes ERROR getting courses for user:", err1)
	}

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"Messages":       msg.GetMessages(c),
			"User":           session.GetLoggedInUserHideError(c),
			"ProfileUser":    profileUser,
			"ProfileCourses": courses,
			"LoggedIn":       session.IsLoggedInValid(c),
			"Meta":           metaDefault,
		},
	)
}

func getNameMedia(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	mediaName := c.Params.ByName("mediaName")
	media, err := db.GetMedia(versionID, mediaName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", media.Type)
	c.Writer.Header().Set("Content-Length", fmt.Sprint(media.Length))

	conn := conn.GetConn()
	WriteMediaChunks(conn, c.Writer, media.ID)
}

func getUserPayouts(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"payout.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}

func getReleaseDelete(c *gin.Context) {
	courseName := c.Params.ByName("course")
	releaseID := c.Query("releaseID")

	release, err := db.GetAllReleaseWithIDStr(releaseID)
	if err != nil {
		log.Println("routes/getReleaseDelete ERROR getting release:", err)
	}

	c.HTML(
		http.StatusOK,
		"confirmDelete.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     session.GetLoggedInUserHideError(c),
			"LoggedIn": session.IsLoggedInValid(c),
			"Meta":     metaDefault,

			// special params for confirmDelete.html
			"Action":  "/" + courseName + "/settings/release/delete/confirm",
			"Message": "Confirm you want to delete release " + fmt.Sprint(release.Num),
			"Data":    release.ID,
			"Further": "This will also delete all versions and user content in this release!",
		},
	)
}
