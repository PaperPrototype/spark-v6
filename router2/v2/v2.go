package v2

import "github.com/gin-gonic/gin"

// every response from the api returns the payload struct
// error messages should be user friendly (since they will get displayed in the front end)
type payload struct {
	Error   string
	Payload interface{}
}

func AddRoutes(group *gin.RouterGroup) {
	group.POST("/course/:courseID", postCourseFORM)                                        // update course info
	group.GET("/course/:courseID/releases", getCourseReleasesJSON)                         // get releases JSON
	group.POST("/course/:courseID/releases", mustBeAuthorCourseID, postCourseReleasesFORM) // create a new release

	group.GET("/releases/:releaseID", getReleaseJSON)                          // get release
	group.POST("/releases/:releaseID", mustBeAuthorReleaseID, postReleaseFORM) // update release
	group.GET("/releases/:releaseID/github", getGithubReleaseJSON)             // get only the github release
	group.GET("/releases/:releaseID/github/tree", mustBeAuthorReleaseID, getGithubReleaseTreeJSON)
	group.GET("/releases/:releaseID/github/tree/assets", mustBeAuthorReleaseID, getGithubReleaseAssetsTreeJSON)
	group.POST("/releases/:releaseID/github", mustBeAuthorReleaseID, postGithubReleaseFORM)   // update or create github release
	group.GET("/releases/:releaseID/sections", getReleaseSectionsJSON)                        // get sections of a release
	group.POST("/releases/:releaseID/section", mustBeAuthorReleaseID, postReleaseSectionFORM) // create a new section
	group.GET("/releases/:releaseID/assets/:name", getReleaseGithubAsset)
	group.GET("/releases/:releaseID/github/resources", getReleaseResourcesJSON)   // get release github resources
	group.GET("/releases/:releaseID/resources/:name", getReleaseResourceDOWNLOAD) // get release github resource
	group.DELETE("/releases/:releaseID", mustBeAuthorReleaseID, deleteRelease)

	// get, edit and delete sections
	group.GET("/sections/:sectionID", getSectionJSON)                                       // get
	group.POST("/sections/:sectionID/github", mustBeAuthorSectionID, postSectionGithubFORM) // create or update github section
	group.POST("/sections/:sectionID", mustBeAuthorSectionID, postSection)                  // edit
	group.DELETE("/sections/:sectionID", mustBeAuthorSectionID, deleteSection)              // delete
	group.GET("/sections/:sectionID/html", getSectionMarkdownHTML)                          // get markdown html

	// must be logged in
	group.GET("/user/github/repos", mustBeLoggedIn, getUserGithubReposJSON)
	group.GET("/user/github/repo/:repoID/branches", mustBeLoggedIn, getUserGithubRepoBranchesJSON)
	group.GET("/user/github/repo/:repoID/branch/:branch/commits", mustBeLoggedIn, getUserGithubRepoBranchCommitsJSON)
}
