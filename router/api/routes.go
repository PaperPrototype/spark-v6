// js api for application frontend
package api

import (
	"main/router/middlewares"

	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup) {
	group.GET("/courses", getCourses)
	group.GET("/section/:sectionID", getSection)
	group.GET("/section/:sectionID/plaintext", mustBeCourseAuthor, getSectionPlaintext)
	group.POST("/section/:sectionID/content/:contentID/edit", postEditSectionContent)
	group.POST("/version/:versionID/posts/new", courseVersionNewPost)
	group.GET("/version/:versionID/posts", getVersionPosts)
	group.GET("/posts/:postID", getPost)
	group.GET("/posts/:postID/plaintext", getPostPlaintext)
	group.POST("/posts/:postID/update", postUpdatePost)

	// logged in users only
	group.GET("/github/user/repos", middlewares.MustBeLoggedIn, getGithubUserRepos)
	group.GET("/github/repo/:repoID/branches", middlewares.MustBeLoggedIn, getGithubRepoBranches)
	group.GET("/github/repo/:repoID/branch/:branch/commits", middlewares.MustBeLoggedIn, getGithubRepoBranchCommits)

	group.GET("/github/users_id/:userID/repos")
	group.GET("/github/users_username/:username/repos")
}
