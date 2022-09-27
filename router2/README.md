## `router.GET("/", getBrowse)`
Get the landing page. When the user is logged in the landing page display's their username.


## `router.NoRoute(getLost)`
404 page not found page. Pretty self explanatory.


## `router.GET("/:username/:course", getCourse)`
Get course. Just like with github, a course is found by the author's username then the course name. This way we don't have problems with name squatting.


## `router.GET("/media/:versionID/name/:mediaName", getMedia)`


This ones special. It gets content (like an image) from the github repo's `Assets` folder, and buffers it through to the front end.

```go
    // readCloser lets us stream the media
    readCloser, err7 := client.Repositories.DownloadContents(ctx, *githubUser.Login, *repo.Name, "Assets/"+mediaName, &github.RepositoryContentGetOptions{
        Ref: githubVersion.SHA,
    })
    if err7 != nil {
        log.Println("routes/get ERROR getting downloading contents in getNameMedia", err6)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    defer readCloser.Close()

    //
    written, err8 := io.Copy(c.Writer, readCloser)
    if err8 != nil {
        log.Println("routes/get ERROR copying/writing contents in getNameMedia", err6)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    c.Writer.Header().Set("Content-Type", mediaType)
    c.Writer.Header().Set("Content-Length", fmt.Sprint(written))
```

This will get deprecated once we move away from "versions" that each hold a git commit (and let us access the repo at that commit), and switch to use single a "release" where *when we update the commit* the version number increases.

# auth

```go
router.POST("/login", postLogin)
router.POST("/signup", postSignup)
router.GET("/logout", getLogout)
```

login, signup, and logout.

```go
router.GET("/login/verify/:verifyUUID", getVerify) // verify account
router.GET("/login/verify/new", getNewVerify)      // send verification email
```

the above are verification for emails, that uses sendgrid to send a confirmation email when you sign up.

# settings

settings is a page that lets the user update their username, update their bio, their full name, and email.

```go
router.GET("/settings", getSettings)
router.POST("/settings/edit/user", middlewares.MustBeLoggedIn, postSettingsEditUser)
router.POST("/settings/edit/email", middlewares.MustBeLoggedIn, postSettingsEditEmail)
```

oauth onboarding flow for github connection
```go
router.GET("/settings/github/connect", middlewares.MustBeLoggedIn, getGithubConnect)
router.GET("/settings/github/connect/return", middlewares.MustBeLoggedIn, getGithubConnectFinished)
```

oauth onboarding flow for stripe connection
```go
router.GET("/settings/stripe/connect", middlewares.MustBeLoggedIn, getStripeConnect)
router.GET("/settings/stripe/login", middlewares.MustBeLoggedIn, getStripeLogin)
router.GET("/settings/stripe/connect/refresh", middlewares.MustBeLoggedIn, getStripeRefresh)
router.GET("/settings/stripe/connect/return", middlewares.MustBeLoggedIn, getStripeConnectFinished)
```