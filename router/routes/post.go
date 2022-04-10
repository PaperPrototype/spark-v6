package routes

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"main/conn"
	"main/db"
	"main/helpers"
	"main/msg"
	"main/payments"
	"main/router/auth"
	"main/upload"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mholt/archiver"
)

func postNew(c *gin.Context) {
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

	// raw name variable
	// lowercase all unique course names
	uncleanName := c.PostForm("name")

	uncleanName = strings.Trim(uncleanName, " ")
	uncleanName = strings.ToLower(uncleanName)

	// name now cleaned
	name := uncleanName

	// fix title
	title := c.PostForm("title")
	title = strings.Trim(title, " ")

	subtitle := c.PostForm("subtitle")

	course := db.Course{
		Name:     name,
		Title:    title,
		Subtitle: subtitle,
		UserID:   user.ID,
	}

	if strings.Contains(uncleanName, " ") {
		c.HTML(
			http.StatusOK,
			"new.html",
			gin.H{
				"Messages": []string{"Spaces are not allowed in course name."},
				"Course":   course,
			},
		)
		return
	}

	if uncleanName == "" || title == "" || subtitle == "" {
		c.HTML(
			http.StatusOK,
			"new.html",
			gin.H{
				"Messages": []string{"All fields must not be empty."},
				"Course":   course,
			},
		)
		return
	}

	available, err1 := db.UserCourseNameAvailable(user.Username, name)
	if !available {
		c.HTML(
			http.StatusOK,
			"new.html",
			gin.H{
				"Messages": []string{"That name is already taken."},
				"Course":   course,
			},
		)
		return
	}

	if err1 != nil {
		log.Println("ERROR:", err1)
		msg.SendMessage(c, "Error. That name is already taken.")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	err := db.CreateCourse(&course)
	if err != nil {
		log.Println("ERROR:", err)
		msg.SendMessage(c, "Error creating new course.")
		c.Redirect(http.StatusFound, "/new")
	}

	// redirect to the new course
	c.Redirect(http.StatusFound, "/"+user.Username+"/"+name+"/settings")
}

func postLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	redirectURL := c.PostForm("redirectURL")

	user, success := db.TryUserPassword(username, password)
	if !success {
		msg.SendMessage(c, "Incorrect username or password.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	sessionToken, err4 := db.CreateSession(user.ID)
	if err4 != nil {
		log.Println("ERROR creating session in db:", err4)
		msg.SendMessage(c, "Error creating session. You will have to login.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	auth.Login(c, sessionToken)

	if redirectURL == "" {
		c.Redirect(http.StatusFound, "/courses")
	}

	c.Redirect(http.StatusFound, redirectURL)
}

func postSignup(c *gin.Context) {
	username := strings.ToLower(c.PostForm("username"))
	name := c.PostForm("name")
	pass := c.PostForm("password")
	confirm := c.PostForm("confirm")
	email := c.PostForm("email")

	available, err4 := db.UsernameAvailable(username)
	if err4 != nil {
		log.Println("routes/post ERROR cheking if username is available:", err4)
		msg.SendMessage(c, "Error checking if username is taken.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if !available {
		msg.SendMessage(c, "Username already taken.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	if strings.Contains(username, " ") {
		msg.SendMessage(c, "Username cannot have spaces.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if !helpers.IsAllowedUsername(username) {
		msg.SendMessage(c, "Invalid username. Allowed characters are "+helpers.AllowedUsernameCharacters)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	emailAvailable, err5 := db.EmailAvailable(email)
	if err5 != nil {
		log.Println("routes/post ERROR cheking if username is available:", err4)
		msg.SendMessage(c, "Error checking if email is taken.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if !emailAvailable {
		msg.SendMessage(c, "That email is already taken.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if pass != confirm {
		msg.SendMessage(c, "Passwords do not match!")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	hash, err1 := helpers.HashPassword(pass)
	if err1 != nil {
		log.Println("ERROR hashing password routes/signup:", err1)
		msg.SendMessage(c, "Password error.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	user := db.User{
		Username: username,
		Name:     name,
		Hash:     hash,
		Email:    email,
	}

	err3 := db.CreateUser(&user)
	if err3 != nil {
		log.Println("ERROR creating user routes/signup:", err3)
		msg.SendMessage(c, "Error creating user.")
		c.Redirect(http.StatusFound, "/signup")
	}

	sessionToken, err4 := db.CreateSession(user.ID)
	if err4 != nil {
		log.Println("ERROR creating session in db:", err4)
		msg.SendMessage(c, "Error creating session. You will have to login.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	auth.Login(c, sessionToken)

	msg.SendMessage(c, "Signed up and logged in!")
	c.Redirect(http.StatusFound, "/courses")
}

func postCourseSettingsDisplay(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	courseID := c.PostForm("courseID")
	title := c.PostForm("title")
	name := c.PostForm("name")
	desc := c.PostForm("desc")
	public := c.PostForm("public")

	available, err := db.UserCourseNameAvailableNotSelf(username, name, courseID)
	if !available {
		msg.SendMessage(c, "That course name is taken.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	if err != nil {
		msg.SendMessage(c, "Error checking if course name available.")
		log.Println("ERROR checking if course name is available:", err)
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	releasesCount := db.CountPublicCourseReleasesLogError(courseID)

	publicBool := false
	// if public and there are releases
	if public != "" && releasesCount != 0 {
		publicBool = true
	}

	if releasesCount == 0 {
		msg.SendMessage(c, "You must have at least 1 public release before you can make a course public.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")

		err1 := db.UpdateCourse(courseID, title, name, desc, false)
		if err1 != nil {
			msg.SendMessage(c, "Error updating course.")
			log.Println("ERROR updating course:", err1)
			c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
			return
		}
		return
	}

	err1 := db.UpdateCourse(courseID, title, name, desc, publicBool)
	if err1 != nil {
		msg.SendMessage(c, "Error updating course.")
		log.Println("ERROR updating course:", err1)
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully updated course!")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postNewRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	desc := c.PostForm("desc")
	price := c.PostForm("price")
	postsNeeded := c.PostForm("postsNeededNum")

	course, err := db.GetUserCoursePreloadUser(username, courseName)
	if err != nil {
		log.Println("routes ERROR getting course:", err)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	priceNumIncorrect, err2 := strconv.ParseUint(price, 10, 64)
	if err2 != nil {
		log.Println("routes ERROR getting course:", err2)
		msg.SendMessage(c, "Error parsing price.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	correctPriceNum := priceNumIncorrect * 100

	if correctPriceNum > payments.MaxCoursePrice {
		msg.SendMessage(c, "The max price of a course is $10 USD")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	author, err3 := db.GetUser(course.UserID)
	if err3 != nil {
		msg.SendMessage(c, "Error getting the course author.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	if !author.HasStripeConnection() {
		correctPriceNum = 0
		msg.SendMessage(c, "You must connect your account to stripe to charge money for a course!")
	}

	postsNeededNum, err5 := strconv.ParseUint(postsNeeded, 10, 16)
	if err5 != nil {
		log.Println("routes ERROR parsing postsNeedeNum for version in postNewGithubVersion:", err5)
		msg.SendMessage(c, "Error parsing number of posts need to complete course version.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	// minimum number of posts is 1
	if postsNeededNum < 1 {
		postsNeededNum = 1
	}

	release := db.Release{
		Num:            course.GetAllNewestCourseReleaseNumLogError() + 1,
		Markdown:       template.HTML(desc),
		CourseID:       course.ID,
		PostsNeededNum: uint16(postsNeededNum),
		Price:          correctPriceNum,
	}

	err1 := db.CreateRelease(&release)
	if err1 != nil {
		log.Println("routes ERROR creating release:", err1)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully created release!")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postNewVersion(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	releaseID := c.PostForm("releaseID")

	fileHandle, err2 := c.FormFile("zipFile")
	if err2 != nil {
		msg.SendMessage(c, "Error getting form file. Make sure you selected a file to upload.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	if filepath.Ext(fileHandle.Filename) != ".zip" {
		msg.SendMessage(c, "Must be a .zip file.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	file, err3 := fileHandle.Open()
	if err3 != nil {
		log.Println("routes ERROR opening form file:", err3)
	}
	defer file.Close()

	uniqueName := uuid.NewString()

	log.Println("routes/post CREATING .zip file.")
	newFile, err4 := os.Create("./" + uniqueName + ".zip")
	if err4 != nil {
		log.Println("routes/post ERROR creating file:", err4)
	}
	defer newFile.Close()
	defer os.Remove("./" + uniqueName + ".zip")

	log.Println("routes/post COPYING zip file")
	numBytesWritten, err5 := io.Copy(newFile, file)
	if err5 != nil {
		log.Println("routes ERROR writing bytes to new file:", err5)
	}
	log.Println("routes/post COPIED file. Bytes written =", numBytesWritten)

	// UNCOMPRESS ZIP FILE

	// the type that will be used to read the input stream
	format := archiver.Zip{}

	// the list of files we want out of the archive; any
	// directories will include all their contents unless
	// we return fs.SkipDir from our handler
	// (leave this nil to walk ALL files from the archive)

	err6 := format.Unarchive("./"+uniqueName+".zip", "./upload"+uniqueName)
	if err6 != nil {
		log.Println("routes ERROR extracting zip file:", err6)
	}
	defer os.RemoveAll("./upload" + uniqueName)

	release, err := db.GetAllRelease(releaseID)
	if err != nil {
		msg.SendMessage(c, "Error getting release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	preview := true
	if c.PostForm("preview") == "" {
		preview = false
	}

	version := db.Version{
		Num:       release.GetNewestVersionNumLogError() + 1,
		ReleaseID: release.ID,
		CourseID:  release.CourseID,
		Preview:   preview,
		Error:     "",
	}

	err1 := db.CreateVersion(&version)
	if err1 != nil {
		log.Println("routes ERROR creating release:", err1)
		msg.SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	courseFolderName := fileHandle.Filename[:len(fileHandle.Filename)-4]

	conn := conn.GetConn()

	err7 := upload.UploadCourse(conn, "./upload"+uniqueName+"/"+courseFolderName, version.ID)
	if err7 != nil {
		msg.SendMessage(c, "Error uploading course: "+fmt.Sprint(err7))
		// log error to version so user can view it and delete version
		upload.LogError(conn, version.ID, err7.Error())
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully created version!")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postEditRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	releaseID := c.PostForm("releaseID")
	desc := c.PostForm("desc")
	price := c.PostForm("price")
	publicStr := c.PostForm("public")
	postsNeeded := c.PostForm("postsNeededNum")

	postsNeededNum, err5 := strconv.ParseUint(postsNeeded, 10, 16)
	if err5 != nil {
		log.Println("routes ERROR parsing postsNeedeNum for version in postNewGithubVersion:", err5)
		msg.SendMessage(c, "Error parsing number of posts need to complete course version.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	// minimum number of posts is 1
	if postsNeededNum < 1 {
		postsNeededNum = 1
	}

	release, err1 := db.GetAllRelease(releaseID)
	if err1 != nil {
		msg.SendMessage(c, "Error updating release.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	priceNumIncorrect, err2 := strconv.ParseUint(price, 10, 64)
	if err2 != nil {
		log.Println("routes ERROR getting course:", err2)
		msg.SendMessage(c, "Error parsing price.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	correctPriceNum := priceNumIncorrect * 100

	if correctPriceNum > payments.MaxCoursePrice {
		msg.SendMessage(c, "THe max price of a course is $10 USD")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	author, err3 := release.GetAuthorUser()
	if err3 != nil {
		msg.SendMessage(c, "Error getting the course author.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	if !author.HasStripeConnection() {
		correctPriceNum = 0
		msg.SendMessage(c, "You must connect your account to stripe to charge money for a course!")
	}

	var public bool = false
	if publicStr != "" && release.HasVersions() {
		public = true
	}

	err := db.UpdateRelease(releaseID, desc, correctPriceNum, public, uint16(postsNeededNum))
	if err != nil {
		log.Println("routes ERROR updating release:", err)
		msg.SendMessage(c, "Error updating release.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	if !release.HasVersions() {
		msg.SendMessage(c, "Release must have uploads before you can make it public.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully updated.")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postDeleteVersion(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	versionID := c.PostForm("versionID")

	err := db.DeleteVersion(versionID)
	if err != nil {
		log.Println("routes ERROR deleting version:", err)
		msg.SendMessage(c, "Error while deleting version")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully deleted version")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postReleaseDeleteConfirm(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")
	releaseID := c.PostForm("data")
	password := c.PostForm("password")

	user, err := auth.GetLoggedInUser(c)
	if err != nil {
		msg.SendMessage(c, "You must be logged in to access that page.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	_, success := db.TryUserPassword(user.Username, password)
	if success {
		err1 := db.DeleteRelease(releaseID)
		if err1 != nil {
			log.Println("routes/postReleaseDeleteConfirm ERROR deleting release:", err1)
			msg.SendMessage(c, "Error. Failed to delete release.")
			c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
			return
		}

		msg.SendMessage(c, "Successfully deleted release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Incorrect password.")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

// create github release or update exisiting github release
func postCreateOrEditGithubRelease(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	releaseID := c.PostForm("releaseID")
	repoID := c.PostForm("repoID")
	branch := c.PostForm("branch")

	if releaseID == "" || repoID == "" || branch == "" {
		msg.SendMessage(c, "Blank fields not allowed.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	log.Println("releaseID:", releaseID, "repoID:", repoID)
	releaseIDNum, err := strconv.ParseUint(releaseID, 10, 64)
	if err != nil {
		log.Println("routes/post ERROR parsing releaseID into uint in postNewGithubRelease:", err)
		msg.SendMessage(c, "Error creating github release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	release, err4 := db.GetAllRelease(releaseIDNum)
	if err4 != nil {
		log.Println("routes/post ERROR getting release in pistGithubRelease:", err4)
		msg.SendMessage(c, "Error getting release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	repoIDNum, err1 := strconv.ParseInt(repoID, 10, 64)
	if err1 != nil {
		log.Println("routes/post ERROR parsing repoID into int in postNewGithubRelease:", err1)
		msg.SendMessage(c, "Error creating github release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	user := auth.GetLoggedInUserLogError(c)
	context := context.Background()
	client := user.GetGithubConnectionLogError().NewClient(context)

	repo, _, err3 := client.Repositories.GetByID(context, repoIDNum)
	if err3 != nil {
		log.Println("routes/post ERROR getting repository by id of user:", err3)
		msg.SendMessage(c, "Error gettign that repository")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	// if has github version, then user cannot change repo, but they can change the branch?
	// create release
	if !release.HasGithubRelease() {
		githubRelease := db.GithubRelease{
			Branch:    branch,
			ReleaseID: releaseIDNum,
			RepoID:    repoIDNum,
			RepoName:  *repo.Name,
		}
		err2 := db.CreateGithubRelease(&githubRelease)
		if err2 != nil {
			log.Println("routes/post ERROR creating githubRelease in postNewGithubRelease:", err2)
			msg.SendMessage(c, "Error creating github release.")
			c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
			return
		}

		msg.SendMessage(c, "Successfully added github repository to release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	githubRelease := release.GetGithubReleaseLogError()
	err2 := db.UpdateGithubRelease(githubRelease.ReleaseID, branch, repoIDNum, *repo.Name)
	if err2 != nil {
		msg.SendMessage(c, "Error updating github course release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully updated github course release.")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}

func postNewGithubVersion(c *gin.Context) {
	username := c.Params.ByName("username")
	courseName := c.Params.ByName("course")

	releaseID := c.PostForm("releaseID")
	sha := c.PostForm("sha")

	if sha == "" {
		msg.SendMessage(c, "You must select a commit.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	if releaseID == "" {
		msg.SendMessage(c, "Error no releaseID.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	release, err2 := db.GetAllRelease(releaseID)
	if err2 != nil {
		log.Println("routes/post ERROR getting release in GetAllRelease:", err2)
		msg.SendMessage(c, "Error getting release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	preview := true
	if c.PostForm("preview") == "" {
		preview = false
	}

	version := db.Version{
		Num:         release.GetNewestVersionNumLogError() + 1,
		ReleaseID:   release.ID,
		CourseID:    release.CourseID,
		Preview:     preview,
		Error:       "",
		UsingGithub: true,
	}
	err1 := db.CreateVersion(&version)
	if err1 != nil {
		log.Println("routes ERROR creating version:", err1)
		msg.SendMessage(c, "Error creating version.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	githubRelease, err4 := db.GetGithubReleaseWithIDStr(releaseID)
	if err4 != nil {
		log.Println("routes/post ERROR getting GithubRelease in postNewGithubVersion:", err4)
		msg.SendMessage(c, "Error getting github release. You may need to update or enable github releases for this release.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	githubVersion := db.GithubVersion{
		VersionID: version.ID,
		RepoID:    githubRelease.RepoID,
		Branch:    githubRelease.Branch,
		SHA:       sha,
	}
	err3 := db.CreateGithubVersion(&githubVersion)
	if err3 != nil {
		log.Println("routes/post ERROR creating githubVersion in postNewGithubVersion:", err3)
		msg.SendMessage(c, "Error creating github version.")
		c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
		return
	}

	msg.SendMessage(c, "Successfully created github version!")
	c.Redirect(http.StatusFound, "/"+username+"/"+courseName+"/settings")
}
