# Release

```go
group.GET("/releases/:releaseID")
group.GET("/releases/:releaseID/github", getGithubReleaseJSON)
group.POST("/releases/:releaseID/github", postGithubReleaseFORM)
group.GET("/releases/:releaseID/sections", getReleaseSectionsJSON)
group.POST("/releases/:releaseID/section", postReleaseSectionFORM)
group.DELETE("/releases/:releaseID/section/:sectionID")
```

# User

getting the repositories of a user
```go
// must be logged in
group.GET("/user/github/repos", mustBeLoggedIn, getUserGithubReposJSON)
```

getting the branches of a repository
```go
group.GET("/user/github/repo/:repoID/branches", mustBeLoggedIn, getUserGithubRepoBranchesJSON)
```

getting the commits of a github repository and branch
```
group.GET("/user/github/repo/:repoID/branch/:branch/commits", mustBeLoggedIn, getUserGithubRepoBranchCommitsJSON)
```