// sparkers interal js api for the website frontend
package api

import (
	"main/router/middlewares"

	"github.com/gin-gonic/gin"
)

func AddRoutes(group *gin.RouterGroup) {
	group.GET("/courses", getCourses)
	group.GET("/level/:level", getLevelCourses)
	group.GET("/version/:versionID/posts", getVersionPosts)                  // portfolio showcase posts
	group.GET("/version/:versionID/posts/showcase", getVersionShowcasePosts) // top rated student work
	group.GET("/version/:versionID/reviews", getCourseReviews)
	group.POST("/version/:versionID/reviews/new", middlewares.MustBeLoggedIn, postNewReview)
	group.POST("/version/:versionID/channel/:channelID/message", middlewares.MustBeLoggedIn, postChannelSendMessage)
	group.GET("/channels/:channelID", getChannelMessages)

	// posts
	group.POST("/posts/new", postNewPost)
	group.GET("/posts/:postID", getPost)
	group.GET("/posts/:postID/comments", getPostComments) // utilizes long polling
	group.POST("/posts/:postID/comment", middlewares.MustBeLoggedIn, postPostComment)
	group.GET("/posts/:postID/plaintext", getPostPlaintext)
	group.POST("/posts/:postID/update", postUpdatePost)

	// notifications
	group.GET("/notifications/newest", middlewares.MustBeLoggedIn, getNewNotifications)
	group.POST("/notifications/done", middlewares.MustBeLoggedIn, postDoneNotification) // set notification as read

	// getting an UPLOAD based course
	group.GET("/section/:sectionID", getSection)
	group.GET("/section/:sectionID/plaintext", mustBeCourseAuthor, getSectionPlaintext)
	group.POST("/section/:sectionID/content/:contentID/edit", postEditSectionContent)

	// getting a GITHUB based course
	// for public viewing and paying customers
	group.GET("/github/version/:versionID/tree", getGithubRepoCommitTree)                         // get the contents in the repo at the commit for that version
	group.GET("/github/version/:versionID/content/:commit_sha/*path", getGithubRepoCommitContent) // used primarily to get an english.md file from the repo

	// get github info for logged in users with their githubConnection
	group.GET("/github/user/repos", middlewares.MustBeLoggedIn, getGithubUserRepos)
	group.GET("/github/repo/:repoID/branches", middlewares.MustBeLoggedIn, getGithubRepoBranches)
	group.GET("/github/repo/:repoID/branch/:branch/commits", middlewares.MustBeLoggedIn, getGithubRepoBranchCommits)
}
