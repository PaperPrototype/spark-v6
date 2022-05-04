package routes

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"main/db"
	"main/mailer"
	"main/markdown"
	"main/msg"
	"main/router/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

var metaDefault = Meta{
	Title: "Sparker - Coding Courses",
	Desc:  "It's time to ditch degree's and switch to portfolio's.",
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

	user, getUserErr := auth.GetLoggedInUser(c)

	// if private and outside user
	if !course.Public && course.UserID != user.ID {
		// deny existence of a course
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
				"Meta": Meta{
					Title: "Sparker - " + course.Title,
					Desc:  course.Subtitle,
				},
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
		if getUserErr == nil {
			if user.HasPurchasedRelease(release.ID) {
				purchased = true
			}
		}
	}

	c.HTML(
		http.StatusOK,
		"course.html",
		gin.H{
			"Version":   release.GetNewestVersionLogError(),
			"Purchased": purchased,
			"Course":    course,
			"Release":   release,
			"Messages":  msg.GetMessages(c),
			"User":      auth.GetLoggedInUserLogError(c),
			"LoggedIn":  auth.IsLoggedInValid(c),
			"Meta": Meta{
				Title: "Sparker - " + course.Title,
				Desc:  course.Subtitle,
			},
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

	channelID := c.Query("channel_id")

	postID := c.Query("post_id")

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
	if !version.HasGithubVersion() {
		section, err2 = version.GetFirstSectionPreload()
		if err2 != nil {
			log.Println("routes/get ERROR getting versions first section in getCourseVersion:", err2)
		}
	}

	var progress int64
	user := auth.GetLoggedInUserLogError(c)
	postsCount := release.UserPostsCountLogError(user.ID)

	if auth.IsLoggedInValid(c) {
		if release.PostsNeededNum != 0 {
			// convert to float for division
			floatProgress := float64(postsCount) / float64(release.PostsNeededNum)

			// convert deciaml to percentage
			floatProgress *= 100

			// cast and save
			progress = int64(floatProgress)
		}
	}

	channels, err4 := db.GetChannels(course.ID)
	if err4 != nil {
		log.Println("routes/get ERROR getting channels in getCourseVersion:", err4)
	}

	channel := &db.Channel{}
	if len(channels) != 0 && channelID == "" {
		channel = &channels[0]
	}

	if channelID != "" {
		var err5 error
		channel, err5 = db.GetChannel(channelID)
		if err5 != nil {
			log.Println("routes/get ERROR getting channel in getCourseVersion:", err5)
			if len(channels) != 0 {
				channel = &channels[0]
			}
		}
	}

	c.HTML(
		http.StatusOK,
		"courseView.html",
		gin.H{
			"ChannelID":  channelID, // if this is not an empty string then we load the channel since it was sent from a notification!
			"Channel":    channel,
			"Channels":   channels,
			"PostID":     postID,
			"PostsCount": postsCount,
			"Release":    release,
			"Course":     course,
			"Version":    version,
			"Section":    section,
			"Messages":   msg.GetMessages(c),
			"User":       auth.GetLoggedInUserLogError(c),
			"LoggedIn":   auth.IsLoggedInValid(c),
			"Meta": Meta{
				Title: "View - " + course.Title,
				Desc:  course.Subtitle,
			},
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

	user, getUserErr := auth.GetLoggedInUser(c)

	// if private and outside user
	if !course.Public && course.UserID != user.ID {
		// deny existence of a course
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
				"Meta": Meta{
					Title: "Sparker - " + course.Title,
					Desc:  course.Subtitle,
				},
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
		if getUserErr == nil {
			if user.HasPurchasedRelease(release.ID) {
				purchased = true
			}
		}
	}

	c.HTML(
		http.StatusOK,
		"course.html",
		gin.H{
			"Version":   release.GetNewestVersionLogError(),
			"Purchased": purchased,
			"Course":    course,
			"Release":   release,
			"Messages":  msg.GetMessages(c),
			"User":      auth.GetLoggedInUserLogError(c),
			"LoggedIn":  auth.IsLoggedInValid(c),
			"Meta": Meta{
				Title: "Sparker - " + course.Title,
				Desc:  course.Subtitle,
			},
		},
	)
}

func getCourseVersionSection(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	postID := c.Query("post_id")

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
	if !version.HasGithubVersion() {
		section, err2 = version.GetFirstSectionPreload()
		if err2 != nil {
			log.Println("routes/get ERROR getting versions first section in getCourseVersion:", err2)
		}
	}

	var progress int64
	user := auth.GetLoggedInUserLogError(c)
	postsCount := release.UserPostsCountLogError(user.ID)

	if auth.IsLoggedInValid(c) {
		if release.PostsNeededNum != 0 {
			// convert to float for division
			floatProgress := float64(postsCount) / float64(release.PostsNeededNum)

			// convert deciaml to percentage
			floatProgress *= 100

			// cast and save
			progress = int64(floatProgress)
		}
	}

	channels, err4 := db.GetChannels(course.ID)
	if err4 != nil {
		log.Println("routes/get ERROR getting channels in getCourseVersion:", err4)
	}

	channel := db.Channel{}
	if len(channels) != 0 {
		channel = channels[0]
	}

	c.HTML(
		http.StatusOK,
		"courseView.html",
		gin.H{
			"Channel":    channel,
			"Channels":   channels,
			"SHA":        c.Params.ByName("sha"),
			"PostID":     postID,
			"PostsCount": postsCount,
			"Release":    release,
			"Course":     course,
			"Version":    version,
			"Section":    section,
			"Messages":   msg.GetMessages(c),
			"User":       auth.GetLoggedInUserLogError(c),
			"LoggedIn":   auth.IsLoggedInValid(c),
			"Meta": Meta{
				Title: "View - " + course.Title,
				Desc:  course.Subtitle,
			},
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

	courses, err1 := profileUser.GetPublicPurchasedCourses()
	if err1 != nil {
		log.Println("routes ERROR getting courses for user:", err1)
	}

	authoredCourses, err2 := profileUser.GetPublicAuthoredCourses()
	if err2 != nil {
		log.Println("routes ERROR getting authored courses for user:", err2)
	}

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"Messages":        msg.GetMessages(c),
			"User":            auth.GetLoggedInUserLogError(c),
			"ProfileUser":     profileUser,
			"ProfileCourses":  courses,
			"AuthoredCourses": authoredCourses,
			"LoggedIn":        auth.IsLoggedInValid(c),
			"Meta":            metaDefault,
		},
	)
}

func getNameMedia(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	mediaName := c.Params.ByName("mediaName")

	version, err1 := db.GetVersion(versionID)
	if err1 != nil {
		log.Println("routes/get ERROR getting version in getNameMedia:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	/* TODO?
	maybe not necessary since course versions can't be viewed unless the course is free (or the user has paid)
	and getting access to the image links without access to the course would be difficult
	*/
	// check if course release is free
	// if paid
	//	 check if student has access to course
	// else
	// 	 free so anyone can view it?

	// if it is a github based version
	if version.HasGithubVersion() {
		githubVersion, err2 := version.GetGithubVersion()
		if err2 != nil {
			log.Println("routes/get ERROR getting githubVersion in getNameMedia:", err2)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err3 := version.GetAuthorUser()
		if err3 != nil {
			log.Println("routes/get ERROR getting githubVersion in getNameMedia:", err2)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		githubConnection, err4 := user.GetGithubConnection()
		if err4 != nil {
			log.Println("routes/get ERROR getting authors github connection in getNameMedia", err4)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		client := githubConnection.NewClient(ctx)

		githubUser, _, err5 := client.Users.Get(ctx, "")
		if err5 != nil {
			log.Println("routes/get ERROR getting githubUser in getNameMedia", err5)
			c.AbortWithStatus(http.StatusNotFound) // user should not know of the existence of this file
			return
		}

		repo, _, err6 := client.Repositories.GetByID(ctx, githubVersion.RepoID)
		if err6 != nil {
			log.Println("routes/get ERROR getting repo by ID in getNameMedia", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		mediaType := filepath.Ext(mediaName)

		if mediaType == ".zip" {
			readCloser, err7 := client.Repositories.DownloadContents(ctx, *githubUser.Login, *repo.Name, "Resources/"+mediaName, &github.RepositoryContentGetOptions{
				Ref: githubVersion.SHA,
			})
			if err7 != nil {
				log.Println("routes/get ERROR getting downloading contents in getNameMedia", err6)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer readCloser.Close()

			written, err8 := io.Copy(c.Writer, readCloser)
			if err8 != nil {
				log.Println("routes/get ERROR copying/writing contents in getNameMedia", err6)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Writer.Header().Set("Content-Type", mediaType)
			c.Writer.Header().Set("Content-Length", fmt.Sprint(written))
			return
		}

		readCloser, err7 := client.Repositories.DownloadContents(ctx, *githubUser.Login, *repo.Name, "Assets/"+mediaName, &github.RepositoryContentGetOptions{
			Ref: githubVersion.SHA,
		})
		if err7 != nil {
			log.Println("routes/get ERROR getting downloading contents in getNameMedia", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer readCloser.Close()

		written, err8 := io.Copy(c.Writer, readCloser)
		if err8 != nil {
			log.Println("routes/get ERROR copying/writing contents in getNameMedia", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Writer.Header().Set("Content-Type", mediaType)
		c.Writer.Header().Set("Content-Length", fmt.Sprint(written))
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
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

func getHome(c *gin.Context) {
	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes ERROR gettingUserWithUsername:", err)
		msg.SendMessage(c, "You must be logged in to view that page.")
		c.Redirect(http.StatusFound, "/courses")
		return
	}

	courses, err1 := user.GetPublicAndPrivatePurchasedCourses()
	if err1 != nil {
		log.Println("routes ERROR getting courses for user:", err1)
	}

	authoredCourses, err2 := user.GetPublicAndPrivateAuthoredCourses()
	if err2 != nil {
		log.Println("routes ERROR getting authored courses for user:", err2)
	}

	c.HTML(
		http.StatusOK,
		"home.html",
		gin.H{
			"Messages":        msg.GetMessages(c),
			"User":            user,
			"ProfileCourses":  courses,
			"AuthoredCourses": authoredCourses,
			"LoggedIn":        auth.IsLoggedInValid(c),
			"Meta":            metaDefault,
		},
	)
}

func getJoin(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"join.html",
		gin.H{
			"Messages": msg.GetMessages(c),
			"User":     auth.GetLoggedInUserLogError(c),
			"LoggedIn": auth.IsLoggedInValid(c),
			"Meta":     metaDefault,
		},
	)
}
