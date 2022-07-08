package api

import (
	"log"
	"main/db"
	"main/router/auth"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func postNewPost(c *gin.Context) {
	if !auth.IsLoggedInValid(c) {
		log.Println("api LOG not logged in valid")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := auth.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api ERROR couldn't get logged in user:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	markdown := c.PostForm("markdown")

	// prevent from posting empty posts
	if strings.Trim(markdown, " ") == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	post := db.Post{
		UserID:   user.ID,
		Markdown: markdown,
	}
	err2 := db.CreatePost(&post)
	if err2 != nil {
		log.Println("api ERROR creating post:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println("posted new post:", post)
	log.Println("markdown is:", markdown)

	versionID := c.Query("version_id")

	// set the user in the post
	post.User = *user

	// reply with the post's markdown
	defer c.JSON(
		http.StatusOK,
		post,
	)

	if versionID == "" {
		// no version_id so return
		return
	}

	log.Println("version_id is present")

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api ERROR getting version:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	/*
		free courses do not need to be "purchased"
		until you want to actually start writing a post for it
		below follows the logic for this
		- if user has purchased the course
			- if course price is zero
				- give free course purchase!
			- else prevent from posting, since they can't have access
		- else prevent from posting, since they can't have access
	*/

	// check if user has purchased the course
	if !db.UserCanAccessCourseRelease(user.ID, version) {
		release, err4 := db.GetPublicRelease(version.ReleaseID)
		if err4 != nil {
			log.Println("api ERROR getting release:", err4)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// if course price is zero, give free course purchase!
		if release.Price == 0 {
			purchase := db.Purchase{
				VersionID:  version.ID,
				UserID:     user.ID,
				ReleaseID:  release.ID,
				CourseID:   release.CourseID,
				AmountPaid: 0,
				AuthorsCut: 0,
				CreatedAt:  time.Now(),
			}

			err5 := db.CreatePurchase(&purchase)
			if err5 != nil {
				log.Println("api ERROR creating purchase:", err5)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("api LOG course is not free")
			// else prevent from posting, since they can't have access
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	postToCourse := db.PostToCourse{
		PostID:    post.ID,
		ReleaseID: version.ReleaseID,
		CourseID:  version.CourseID,
	}
	err3 := db.CreatePostToCourse(&postToCourse)
	if err3 != nil {
		log.Println("api ERROR creating postToRelease:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func postUpdatePost(c *gin.Context) {
	if !auth.IsLoggedInValid(c) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := auth.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api ERROR couldn't get logged in user:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	postID := c.Params.ByName("postID")

	post, err := db.GetPostPreloadUser(postID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if post.UserID != user.ID {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newMarkdown := c.PostForm("markdown")

	if strings.Trim(newMarkdown, " ") == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	db.UpdatePost(postID, newMarkdown)
}

func postEditSectionContent(c *gin.Context) {
	sectionID := c.Params.ByName("sectionID")
	contentID := c.Params.ByName("contentID")

	contentMarkdown := c.PostForm("content")
	versionID := c.PostForm("versionID")

	if !auth.IsLoggedInValid(c) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err1 := auth.GetLoggedInUser(c)
	if err1 != nil {
		log.Println("api ERROR couldn't get logged in user:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	version, err2 := db.GetVersion(versionID)
	if err2 != nil {
		log.Println("api ERROR getting version:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	release, err4 := db.GetAllRelease(version.ReleaseID)
	if err4 != nil {
		log.Println("api ERROR getting release:", err4)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	course, err3 := db.GetCourseWithIDPreloadUser(release.CourseID)
	if err3 != nil {
		log.Println("api ERROR getting course:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user.ID != course.UserID {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err := db.UpdateSectionContentAndIncreasePatch(sectionID, contentID, contentMarkdown, versionID)
	if err != nil {
		log.Println("api ERROR updating section content for api/postEditSectionContent:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func postNewReview(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	markdown := c.PostForm("markdown")
	ratingStr := c.PostForm("rating")

	if markdown == "" || ratingStr == "" {
		// cannot post an empty post
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	version, err := db.GetVersion(versionID)
	if err != nil {
		log.Println("api/post ERROR getting version in postNewReview:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user := auth.GetLoggedInUserLogError(c)

	numberOfReviews := db.CountUserReviewsLogError(user.ID, version.CourseID)
	if numberOfReviews >= 1 {
		// user can only post 1 review per course
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	post := db.Post{
		UserID:   user.ID,
		Markdown: markdown,
	}
	err1 := db.CreatePost(&post)
	if err1 != nil {
		log.Println("api/post ERROR creating Post in postNewReview:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rating, err2 := strconv.ParseUint(ratingStr, 10, 8)
	if err2 != nil {
		log.Println("api/post ERROR parsing rating num in postNewReview:", err2)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// keep in range of 0 to 5
	if rating > 5 {
		rating = 5
	}

	review := db.PostToCourseReview{
		CourseID:  version.CourseID,
		ReleaseID: version.ReleaseID,
		PostID:    post.ID,
		Rating:    uint8(rating),
		UserID:    user.ID,
	}
	err3 := db.CreateReview(&review)
	if err3 != nil {
		log.Println("api/post ERROR creating post PostToCourseReview in postNewReview:", err3)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// gin will set the status to ok
}
