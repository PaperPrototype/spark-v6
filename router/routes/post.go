package routes

import (
	"fmt"
	"io"
	"log"
	"main/conn"
	"main/db"
	"main/helpers"
	"main/router/session"
	"main/upload"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mholt/archiver"
)

func postNew(c *gin.Context) {
	if !session.IsLoggedInValid(c) {
		SendMessage(c, "You must be logged in to create a new course.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	user, err2 := session.GetLoggedInUser(c)
	if err2 != nil {
		log.Println("ERROR getting logged in user:", err2)
		SendMessage(c, "Error getting logged in user.")
		c.Redirect(http.StatusFound, "/")
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

	desc := c.PostForm("desc")

	course := db.Course{
		Name:   name,
		Title:  title,
		Desc:   desc,
		UserID: user.ID,
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

	if uncleanName == "" || title == "" || desc == "" {
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

	available, err1 := db.CourseNameAvailable(name)
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
		SendMessage(c, "Error. That name is already taken.")
		c.Redirect(http.StatusFound, "/new")
		return
	}

	err := db.CreateCourse(&course)
	if err != nil {
		log.Println("ERROR:", err)
		SendMessage(c, "Error creating new course.")
		c.Redirect(http.StatusFound, "/new")
	}

	// redirect to the new course
	c.Redirect(http.StatusFound, "/"+name+"/settings")
}

func postLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, success := db.TryGetUser(username, password)
	if !success {
		SendMessage(c, "Incorrect username or password.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	sessionToken, err4 := db.CreateSession(user.ID)
	if err4 != nil {
		log.Println("ERROR creating session in db:", err4)
		SendMessage(c, "Error creating session. You will have to login.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session.Login(c, sessionToken)

	c.Redirect(http.StatusFound, "/")
}

func postSignup(c *gin.Context) {
	username := strings.ToLower(c.PostForm("username"))
	name := c.PostForm("name")
	pass := c.PostForm("password")
	confirm := c.PostForm("confirm")
	email := c.PostForm("email")

	available, err := db.UsernameAvailable(username)
	if !available {
		SendMessage(c, "Username already taken.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if err != nil {
		SendMessage(c, "Username taken or error.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if strings.Contains(username, " ") {
		SendMessage(c, "Username cannot have spaces.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	emailAvailable, err2 := db.EmailAvailable(email)
	if !emailAvailable {
		SendMessage(c, "That email is already taken.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if err2 != nil {
		SendMessage(c, "Email taken or error.")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	if pass != confirm {
		SendMessage(c, "Passwords do not match!")
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	hash, err1 := helpers.HashPassword(pass)
	if err1 != nil {
		log.Println("ERROR hashing password routes/signup:", err1)
		SendMessage(c, "Password error.")
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
		SendMessage(c, "Error creating user.")
		c.Redirect(http.StatusFound, "/signup")
	}

	sessionToken, err4 := db.CreateSession(user.ID)
	if err4 != nil {
		log.Println("ERROR creating session in db:", err4)
		SendMessage(c, "Error creating session. You will have to login.")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session.Login(c, sessionToken)

	SendMessage(c, "Signed up and logged in!")
	c.Redirect(http.StatusFound, "/courses")
}

func postCourseSettingsDisplay(c *gin.Context) {
	courseID := c.PostForm("courseID")
	title := c.PostForm("title")
	name := c.PostForm("name")
	desc := c.PostForm("desc")

	course, err2 := db.GetCourseWithIDStr(courseID)
	if err2 != nil {
		SendMessage(c, "Error finding course.")
		log.Println("ERROR finding course:", err2)
		c.Redirect(http.StatusFound, "/")
		return
	}

	available, err := db.CourseNameAvailableNotSelf(name, courseID)
	if !available {
		SendMessage(c, "That course name is taken.")
		c.Redirect(http.StatusFound, "/"+course.Name+"/settings")
		return
	}

	if err != nil {
		SendMessage(c, "Error checking if course name available.")
		log.Println("ERROR checking if course name is available:", err)
		c.Redirect(http.StatusFound, "/"+course.Name+"/settings")
		return
	}

	err1 := db.UpdateCourse(courseID, title, name, desc)
	if err1 != nil {
		SendMessage(c, "Error updating course.")
		log.Println("ERROR updating course:", err1)
		c.Redirect(http.StatusFound, "/"+course.Name+"/settings")
		return
	}

	SendMessage(c, "Successfully updated course!")
	c.Redirect(http.StatusFound, "/"+name+"/settings")
}

func postNewRelease(c *gin.Context) {
	courseName := c.Params.ByName("course")
	desc := c.PostForm("desc")

	course, err := db.GetCourse(courseName)
	if err != nil {
		log.Println("routes ERROR getting course:", err)
		SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	release := db.Release{
		Num:      course.GetNewestCourseReleaseNumLogError() + 1,
		Desc:     desc,
		CourseID: course.ID,
	}

	err1 := db.CreateRelease(&release)
	if err1 != nil {
		log.Println("routes ERROR creating release:", err1)
		SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	SendMessage(c, "Successfully created release!")
	c.Redirect(http.StatusFound, "/"+courseName+"/settings")
}

func postNewVersion(c *gin.Context) {
	courseName := c.Params.ByName("course")
	releaseID := c.PostForm("releaseID")

	fileHandle, err2 := c.FormFile("zipFile")
	if err2 != nil {
		SendMessage(c, "Error getting form file. Make sure you selected a file to upload.")
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	if filepath.Ext(fileHandle.Filename) != ".zip" {
		SendMessage(c, "Must be a .zip file.")
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	file, err3 := fileHandle.Open()
	if err3 != nil {
		log.Println("routes ERROR opening form file:", err3)
	}
	defer file.Close()

	uniqueName := uuid.NewString()

	newFile, err4 := os.Create("./" + uniqueName + ".zip")
	if err4 != nil {
		log.Println("routes ERROR creating file:", err4)
	}
	defer newFile.Close()
	defer os.Remove("./" + uniqueName + ".zip")

	numBytesWritten, err5 := io.Copy(newFile, file)
	if err5 != nil {
		log.Println("routes ERROR writing bytes to new file:", err5)
	}
	log.Println("file downloaded. Bytes written =", numBytesWritten)

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

	release, err := db.GetReleaseWithIDStr(releaseID)
	if err != nil {
		SendMessage(c, "Error getting release.")
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	version := db.Version{
		Num:       release.GetNewestVersionNumLogError() + 1,
		ReleaseID: release.ID,
		CourseID:  release.CourseID,
	}

	err1 := db.CreateVersion(&version)
	if err1 != nil {
		log.Println("routes ERROR creating release:", err1)
		SendMessage(c, "Error getting course.")
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	courseFolderName := fileHandle.Filename[:len(fileHandle.Filename)-4]

	conn := conn.GetConn()

	err7 := upload.UploadCourse(conn, "./upload"+uniqueName+"/"+courseFolderName, version.ID)
	if err7 != nil {
		SendMessage(c, "Error uploading course: "+fmt.Sprint(err7))
		c.Redirect(http.StatusFound, "/"+courseName+"/settings")
		return
	}

	SendMessage(c, "Successfully created version!")
	c.Redirect(http.StatusFound, "/"+courseName+"/settings")
}

func postEditRelease(c *gin.Context) {
	course := c.Params.ByName("course")
	releaseID := c.PostForm("releaseID")
	desc := c.PostForm("desc")

	err := db.UpdateRelease(releaseID, desc)
	if err != nil {
		log.Println("routes ERROR updating release:", err)
		SendMessage(c, "Error updating release.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	SendMessage(c, "Successfully updated.")
	c.Redirect(http.StatusFound, "/"+course+"/settings")
}

func postDeleteVersion(c *gin.Context) {
	course := c.Params.ByName("course")
	versionID := c.PostForm("versionID")

	err := db.DeleteVersion(versionID)
	if err != nil {
		log.Println("routes ERROR deleting version:", err)
		SendMessage(c, "Error while deleting version")
		c.Redirect(http.StatusFound, "/"+course+"/settings")
		return
	}

	SendMessage(c, "Successfully deleted version")
	c.Redirect(http.StatusFound, "/"+course+"/settings")
}
