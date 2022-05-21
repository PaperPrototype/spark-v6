package api

import (
	"log"
	"main/db"
	"main/markdown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getCourses(c *gin.Context) {
	search := c.Query("search")

	var courses []db.Course
	var err error

	if search == "" {
		courses, err = db.GetAllPublicCoursesPreload()
	} else {
		log.Println("applying search...")
		courses, err = db.GetAllPublicCoursesPreloadAndSearch(search)
	}

	if err != nil {
		log.Println("api ERROR getting course for api/getCourses:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, courses)
}

func getSection(c *gin.Context) {
	sectionID := c.Params.ByName("sectionID")
	section, err := db.GetSectionPreloadConvertMarkdown(sectionID)

	if err != nil {
		log.Println("api ERROR getting section for api/getSection:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, section)
}

func getSectionPlaintext(c *gin.Context) {
	sectionID := c.Params.ByName("sectionID")
	section, err := db.GetSectionPreload(sectionID)

	if err != nil {
		log.Println("api ERROR getting section for api/getSectionPlaintext:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, section)
}

func getVersionPosts(c *gin.Context) {
	versionID := c.Params.ByName("versionID")

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api ERROR getting version for api/getVersionPosts:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	releasePosts, err1 := db.GetReleasePosts(version.ReleaseID, version.CourseID)
	if err1 != nil {
		log.Println("api ERROR getting posts for api/getVersionPosts:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		releasePosts,
	)
}

func getVersionShowcasePosts(c *gin.Context) {
	versionID := c.Params.ByName("versionID")

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api ERROR getting version for api/getVersionPosts:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	releasePosts, err1 := db.GetReleasePostsOrderByLikes(version.ReleaseID, version.CourseID)
	if err1 != nil {
		log.Println("api ERROR getting posts for api/getVersionPosts:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		releasePosts,
	)
}

func getPost(c *gin.Context) {
	postID := c.Params.ByName("postID")
	post, err := db.GetPostPreloadUser(postID)
	if err != nil {
		log.Println("api ERROR getting post for api/getPost:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	buf, err1 := markdown.Convert([]byte(post.Markdown))
	if err1 != nil { // if error
		log.Println("api ERROR converting markdown:", err1)
	} else { // no error
		post.Markdown = buf.String()
	}

	c.JSON(
		http.StatusOK,
		post,
	)
}

func getPostPlaintext(c *gin.Context) {
	postID := c.Params.ByName("postID")
	post, err := db.GetPostPreloadUser(postID)
	if err != nil {
		log.Println("api ERROR getting post for api/getPost:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(
		http.StatusOK,
		post,
	)
}

func getCourseReviews(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api/get ERROR getting course version in getCourseReviews:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	limit, err2 := strconv.ParseInt(limitStr, 10, 64)
	if err2 != nil {
		log.Println("api/get ERROR parsing limitStr in getCourseReviews:", err2)
	}

	offset, err3 := strconv.ParseInt(offsetStr, 10, 64)
	if err2 != nil {
		log.Println("api/get ERROR parsing limitStr in getCourseReviews:", err3)
	}

	reviews, err1 := db.GetCourseReviews(version.CourseID, int(offset), int(limit))
	if err1 != nil {
		log.Println("api/get ERROR getting course reviews in getCourseReviews:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// convert markdown
	for i := range reviews {
		convertedMarkdown, err4 := markdown.Convert([]byte(reviews[i].Post.Markdown))
		if err4 != nil {
			log.Println("api/get ERROR converting markdown in getCourseReviews:", err4)
			continue // skip to next post
		}

		reviews[i].Post.Markdown = convertedMarkdown.String()
	}

	c.JSON(
		http.StatusOK,
		struct {
			Count   int64
			Reviews []db.PostToCourseReview
		}{
			Count:   db.CountCourseReviewsLogError(version.CourseID),
			Reviews: reviews,
		},
	)
}
