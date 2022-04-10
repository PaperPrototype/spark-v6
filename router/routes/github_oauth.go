package routes

import (
	"context"
	"fmt"
	"log"
	"main/db"
	"main/githubapi"
	"main/helpers"
	"main/msg"
	"main/router/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// TODO refresh github oauth token? NOPE! github oauth tokens currently NEVER expire

// You must register the app at https://github.com/settings/applications
// Set callback to http://127.0.0.1:7000/github_oauth_cb
// Set ClientId and ClientSecret to
var (
	oauthConfig = &oauth2.Config{
		ClientID:     helpers.GetGithubClientID(),
		ClientSecret: helpers.GetGithubClientSecret(),
		RedirectURL:  helpers.GetHost() + "/settings/github/connect/return",

		// select level of access you want https://developer.github.com/v3/oauth/#scopes
		Scopes:   []string{"user:email", "repo"},
		Endpoint: githuboauth.Endpoint,
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString string = "jehwgjkbn3qeyi23y98oihnabieyfgh09weohg"
)

func getGithubConnect(c *gin.Context) {
	user := auth.GetLoggedInUserLogError(c)
	if user.HasGithubConnection() {
		msg.SendMessage(c, "Your account is already connected to github")
		c.Redirect(http.StatusFound, "/settings/teaching")
		return
	}

	url := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

func getGithubConnectFinished(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		fmt.Printf("routes/github ERROR invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		msg.SendMessage(c, "Error connecting github account.")
		c.Redirect(http.StatusTemporaryRedirect, "/settings/teaching")
		return
	}

	code := c.Query("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("routes/github ERROR oauthConf.Exchange() failed with '%s'\n", err)
		msg.SendMessage(c, "Error connecting github account.")
		c.Redirect(http.StatusTemporaryRedirect, "/settings/teaching")
		return
	}

	loggedInUser := auth.GetLoggedInUserLogError(c)
	githubConnection := githubapi.GithubConnection{
		UserID:      loggedInUser.ID,
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
	}
	err1 := db.CreateGithubConnection(&githubConnection)
	if err1 != nil {
		log.Println("routes/github ERROR creating github connection in getGithubConnectFinished:", err1)
		msg.SendMessage(c, "Error connecting github account.")
		c.Redirect(http.StatusTemporaryRedirect, "/settings/teaching")
		return
	}

	oauthClient := oauthConfig.Client(context.Background(), token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Printf("routes/github ERROR client.Users.Get() failed with '%s'\n", err)
		msg.SendMessage(c, "Error connecting github account.")
		c.Redirect(http.StatusTemporaryRedirect, "/settings/teaching")
		return
	}

	fmt.Printf("routes/github SUCCESS Logged in as GitHub user: %s\n", *user.Login)
	msg.SendMessage(c, "Successfully connected github account!")
	c.Redirect(http.StatusFound, "/settings")
}
