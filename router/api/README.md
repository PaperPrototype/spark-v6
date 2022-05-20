# Docs for API package

## posts routes

```go
	// in /router/router.go
	api.AddRoutes(router.Group("/api"))

	// in /router/api/routes.go
	group.POST("/posts/new", postNewPost)
	group.GET("/posts/:postID", getPost)
	group.GET("/posts/:postID/comments", getPostComments) // utilizes long polling
	group.POST("/posts/:postID/comment", middlewares.MustBeLoggedIn, postPostComment)
	group.GET("/posts/:postID/plaintext", getPostPlaintext)
	group.POST("/posts/:postID/update", postUpdatePost)
```

### /api/posts/new
params:
	- "markdown" (FORM)
	- "version_id" (URL QUERY)

description:
	- creates a new post and returns the created post with the .User field pre-loaded
	- if the "version_id" query parameter shows up then the post will be added to a course release