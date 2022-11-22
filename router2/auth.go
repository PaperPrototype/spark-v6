package router2

import (
	"log"
	"main/auth2"
	"main/db"
	"main/helpers"
	"main/mailer"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func postLogin(c *gin.Context) {
	email := c.PostForm("username") // form field name has to be "username" for autcomplete
	password := c.PostForm("password")
	redirectURL := c.PostForm("redirectURL")

	user, success := db.TryEmailPassword(email, password)
	if !success {
		msg.SendMessage(c, "Incorrect email or password.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	sessionToken, err4 := db.CreateSession(user.ID)
	if err4 != nil {
		log.Println("ERROR creating session in db:", err4)
		msg.SendMessage(c, "Error creating session. You will have to login.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	auth2.Login(c, sessionToken)

	if redirectURL == "" {
		c.Redirect(http.StatusFound, redirectURL)
	}

	c.Redirect(http.StatusFound, redirectURL)
}

func postSignup(c *gin.Context) {
	pass := c.PostForm("password")
	confirm := c.PostForm("confirm")
	email := c.PostForm("username") // form name must be "username" for autocomplete
	name := c.PostForm("name")
	redirectURL := c.PostForm("redirectURL")

	if email == "" {
		msg.SendMessage(c, "Error. Blank email not allowed.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if pass == "" || confirm == "" {
		msg.SendMessage(c, "Error. Blank password not allowed.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	emailAvailable, err5 := db.EmailAvailable(email)
	if err5 != nil {
		log.Println("routes/post ERROR checking if email is available:", err5)
		msg.SendMessage(c, "Error checking if email is taken.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if !emailAvailable {
		msg.SendMessage(c, "That email is already taken.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if pass != confirm {
		msg.SendMessage(c, "Passwords do not match!")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	hash, err1 := helpers.HashPassword(pass)
	if err1 != nil {
		log.Println("ERROR hashing password routes/signup:", err1)
		msg.SendMessage(c, "Password error.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	user := db.User{
		Name:     name,
		Username: "default_username_" + uuid.NewString(),
		Hash:     hash,
		Email:    email,
	}
	err3 := db.CreateUser(&user)
	if err3 != nil {
		log.Println("ERROR creating user routes/signup:", err3)
		msg.SendMessage(c, "Error creating user.")
		c.Redirect(http.StatusFound, redirectURL)
	}

	sessionToken, err4 := db.CreateSession(user.ID)
	if err4 != nil {
		log.Println("ERROR creating session in db:", err4)
		msg.SendMessage(c, "Error creating session. You will have to login.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	auth2.Login(c, sessionToken)

	msg.SendMessage(c, "Sign up successful. Check your email to verify your account. Make sure to check your spam folder.")

	err6 := mailer.SendVerification(user.ID)
	if err6 != nil {
		log.Println("router2/auth.go ERROR sending email verifcation in postSignup:", err6)
		msg.SendMessage(c, "Account creating successful. Error sending email verification.")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	c.Redirect(http.StatusFound, redirectURL)
}

func getLogout(c *gin.Context) {
	auth2.Logout(c)
	c.Redirect(http.StatusFound, "/")
}

func getNewVerify(c *gin.Context) {
	user, err := auth2.GetLoggedInUser(c)
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

	msg.SendMessage(c, message)
	c.Redirect(http.StatusFound, "/settings")
}
