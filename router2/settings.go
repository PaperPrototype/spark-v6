package router2

import (
	"log"
	"main/auth2"
	"main/db"
	"main/helpers"
	"main/msg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getSettings(c *gin.Context) {
	user, err := auth2.GetLoggedInUser(c)
	if err != nil {
		log.Println("router2/settings.go ERROR getting user:", err)
		msg.SendMessage(c, "Failed to get logged in user")
		c.Redirect(http.StatusFound, "/")
		return
	}

	stripeConnection, err1 := db.GetStripeConnection(user.ID)
	if err1 != nil {
		log.Println("routes/settings ERROR getting stripeConnection in getSettingsCourses:", err1)
	}

	c.HTML(
		http.StatusOK,
		"settings_.html",
		gin.H{
			"StripeConnection": stripeConnection,
			"PayoutsEnabled":   stripeConnection.PayoutsEnabledLogError(),
			"User":             auth2.GetLoggedInUserLogError(c),
			"LoggedIn":         auth2.IsLoggedInValid(c),
			"Messages":         msg.GetMessages(c),
		},
	)
}

func postSettingsEditUser(c *gin.Context) {
	user, err := auth2.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes/settings ERROR getting user in postSettingsEditUser:", err)
		msg.SendMessage(c, "Error getting logged in user.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	username := c.PostForm("username")
	name := c.PostForm("name")
	bio := c.PostForm("bio")

	log.Println("bio:", bio)

	// if username changed
	if username != user.Username {
		if !helpers.IsAllowedUsername(username) {
			if !db.UsernameAvailableIgnoreError(username) {
				err2 := db.UpdateUser(user.ID, user.Username, name, bio, user.Email)
				if err2 != nil {
					log.Println("routes/settings ERROR updating user in postSettingsEditUser:", err2)
					msg.SendMessage(c, "Error updating user.")
					c.Redirect(http.StatusFound, "/settings")
					return
				}

				msg.SendMessage(c, "That username is already taken.")
				c.Redirect(http.StatusFound, "/settings")
				return
			}

			msg.SendMessage(c, "Username contained invalid characters. Valid characters are "+helpers.AllowedUsernameCharacters)
			username = helpers.ConvertToAllowedName(username)
		}
	}

	err1 := db.UpdateUser(user.ID, username, name, bio, user.Email)
	if err1 != nil {
		log.Println("routes/settings ERROR updating user in postSettingsEditUser:", err1)
		msg.SendMessage(c, "Error updating username.")
		c.Redirect(http.StatusFound, "/settings")
		return
	}

	msg.SendMessage(c, "Successfully updated user.")
	c.Redirect(http.StatusFound, "/settings")
}

func postSettingsEditEmail(c *gin.Context) {
	user, err := auth2.GetLoggedInUser(c)
	if err != nil {
		log.Println("routes/settings ERROR getting user in postSettingsEditUser:", err)
		msg.SendMessage(c, "Error getting logged in user.")
		c.Redirect(http.StatusFound, "/")
		return
	}

	email := c.PostForm("email")

	if user.Email == email {
		msg.SendMessage(c, "Email is the same. Nothing to update.")
		c.Redirect(http.StatusFound, "/settings")
		return
	} else {
		if !db.EmailAvailableIgnoreError(email) {
			msg.SendMessage(c, "That email is taken.")
			c.Redirect(http.StatusFound, "/settings")
			return
		}

		err2 := user.SetVerified(false)
		if err2 != nil {
			log.Println("routes/settings ERROR updating user in postSettingsEditUser:", err2)
			msg.SendMessage(c, "Error updating user's email.")
			c.Redirect(http.StatusFound, "/settings")
			return
		}

		err1 := db.UpdateUser(user.ID, user.Username, user.Name, user.Bio, email)
		if err1 != nil {
			log.Println("routes/settings ERROR updating user in postSettingsEditUser:", err1)
			msg.SendMessage(c, "Error updating user's email.")
			c.Redirect(http.StatusFound, "/settings")
			return
		}
	}

	msg.SendMessage(c, "Successfully updated email.")
	c.Redirect(http.StatusFound, "/settings")
}
