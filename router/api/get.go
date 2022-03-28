package api

import (
	"log"
	"main/db"
	"main/markdown"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCourses(c *gin.Context) {
	courses, err := db.GetAllCourses()

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

	releasePosts, err1 := db.GetReleasePosts(version.ReleaseID)
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
