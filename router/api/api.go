// js api for application frontend
package api

import (
	"main/router/middlewares"

	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup) {
	group.GET("/courses", getCourses)
	group.GET("/section/:sectionID", getSection)
	group.GET("/section/:sectionID/plaintext", getSectionPlaintext)
	group.POST("/section/:sectionID/content/:contentID/edit", postEditSectionContent)
	group.POST("/version/:versionID/posts/new", courseVersionNewPost)
	group.GET("/version/:versionID/posts", getVersionPosts)
	group.GET("/posts/:postID", getPost)
	group.GET("/posts/:postID/plaintext", getPostPlaintext)
	group.POST("/posts/:postID/update", postUpdatePost)

	// logged in users only
	group.GET("/github/user/repos", middlewares.MustBeLoggedIn, getGithubUserRepos)

	// not logged in users
	group.GET("/github/users/:userID/repos")
}
