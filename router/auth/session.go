package auth

import (
	"log"
	"main/db"
	"strings"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, sessionToken string) {
	c.SetCookie("session", sessionToken, 0, "/", c.Request.URL.Hostname(), true, true)
}

// returns false if no session cookie
// returns false if session cookie is invalid
func IsLoggedInValid(c *gin.Context) bool {
	// delete expired sessions
	err := db.DeleteExpiredSessions()
	if err != nil {
		log.Println("session ERROR deleting old session? (possilbe there is no sessions):", err)
	}

	// get cookie session
	cookie, _ := c.Cookie("session")

	// get rid of empty space and check if cookie is empty
	if strings.Trim(cookie, " ") == "" {
		return false
	}

	// if session actually exists and is valid
	if !db.SessionExists(cookie) {
		return false
	}

	return true
}

func Logout(c *gin.Context) {
	// get cookie session
	cookie, err := c.Cookie("session")

	// if no error
	if err != nil {
		err1 := db.DeleteSession(cookie)
		if err1 != nil {
			log.Println("session ERROR deletinng session:", err1)
		}
	}

	c.SetCookie("session", "", -2, "/", c.Request.URL.Hostname(), true, true)
}

func GetSessionToken(c *gin.Context) string {
	token, _ := c.Cookie("session")
	return token
}

func GetLoggedInUser(c *gin.Context) (*db.User, error) {
	// delete expired sessions
	err := db.DeleteExpiredSessions()
	if err != nil {
		log.Println("session ERROR deleting old session? (possilbe there is no sessions):", err)
	}

	user, err := db.GetUserFromSession(GetSessionToken(c))
	if err != nil {
		log.Println("ERROR finding user for that session:", err)
	}

	return user, err
}

func GetLoggedInUserLogError(c *gin.Context) *db.User {
	// delete expired sessions
	err := db.DeleteExpiredSessions()
	if err != nil {
		log.Println("session ERROR deleting old session? (possilbe there is no sessions):", err)
	}

	user, err1 := db.GetUserFromSession(GetSessionToken(c))
	if err1 != nil {
		log.Println("ERROR finding user for that session:", err1)
	}

	return user
}
