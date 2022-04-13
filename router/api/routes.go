// js api for application frontend
package api

import (
	"main/router/middlewares"

	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup) {
	group.GET("/courses", getCourses)
	group.GET("/version/:versionID/posts/portfolio", getVersionPortfolioPosts)                           // portfolio proof of work posts
	group.GET("/version/:versionID/posts/proposal", getVersionProposalPosts)                             // proposal posts for final project
	group.GET("/version/:versionID/projects", getVersionProjects)                                        // course final projects
	group.POST("/version/:versionID/posts/:postID/comment", middlewares.MustBeLoggedIn, postPostComment) // creates notification linked to course release

	group.GET("/posts/:postID", getPost)
	group.GET("/posts/:postID/comments", getPostComments) // utilizes long polling
	group.GET("/posts/:postID/plaintext", getPostPlaintext)
	group.POST("/posts/:postID/update", postUpdatePost)

	// notifications
	group.GET("/notifications/newest", middlewares.MustBeLoggedIn, getNewNotifications)
	group.POST("/notifications/done", middlewares.MustBeLoggedIn, postDoneNotification) // set notification as read
	group.GET("/notifications", middlewares.MustBeLoggedIn)                             // get all notifications in pages

	// getting an UPLOAD based course
	group.GET("/section/:sectionID", getSection)
	group.GET("/section/:sectionID/plaintext", mustBeCourseAuthor, getSectionPlaintext)
	group.POST("/section/:sectionID/content/:contentID/edit", postEditSectionContent)
	group.POST("/version/:versionID/posts/new", courseVersionNewPost)

	// getting a GITHUB based course
	// for public viewing and paying customers
	group.GET("/github/version/:versionID/tree", getGithubRepoCommitTree)
	group.GET("/github/version/:versionID/content/:commit_sha/*path", getGithubRepoCommitContent)

	// uses logged in users github connection
	group.GET("/github/user/repos", middlewares.MustBeLoggedIn, getGithubUserRepos)
	group.GET("/github/repo/:repoID/branches", middlewares.MustBeLoggedIn, getGithubRepoBranches)
	group.GET("/github/repo/:repoID/branch/:branch/commits", middlewares.MustBeLoggedIn, getGithubRepoBranchCommits)
}
