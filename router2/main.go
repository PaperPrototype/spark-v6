package router2

import (
	"html/template"
	"log"
	"main/auth2"
	"main/db"
	"main/markdown"
	"main/msg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getBrowse(c *gin.Context) {
	courses, _ := db.GetAllPublicCoursesPreload()

	ownedCourses := []db.Ownership{}
	user, err := auth2.GetLoggedInUser(c)
	if err == nil {
		ownedCourses, _ = db.GetOwnershipsPreloadCourses(user.ID)
	}

	authoredCourses := []db.Course{}
	if err == nil {
		authoredCourses, _ = user.GetPublicAndPrivateAuthoredCourses()
	}

	c.HTML(
		http.StatusOK,
		"browse_.html",
		gin.H{
			"AuthoredCourses": authoredCourses,
			"OwnedCourses":    ownedCourses,
			"User":            user,
			"LoggedIn":        auth2.IsLoggedInValid(c),
			"Messages":        msg.GetMessages(c),
			"Courses":         courses,
			"Meta": meta{
				Title:    "Sparker - Browse",
				Desc:     "Learn coding to build ideas",
				ImageURL: "/resources2/images/sparker_code_hl_banner.png",
			},
		},
	)
}

func getCourse(c *gin.Context) {
	usernameParam := c.Params.ByName("username")
	courseParam := c.Params.ByName("course")
	sectionIDParam := c.Params.ByName("sectionID")
	_ = c.Params.ByName("releaseID")

	course, err := db.GetUserCoursePreload(usernameParam, courseParam)
	if err != nil {
		log.Println("router/get.go ERROR getting course:", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// if author then they can view private releases
	releases := []db.Release{}
	if auth2.GetLoggedInUserLogError(c).ID == course.UserID {
		releases, _ = db.GetAnyReleases(course.ID)
	} else {

		releases, _ = db.GetPublicReleases(course.ID)
	}

	owned := false
	user, err1 := auth2.GetLoggedInUser(c)
	if err1 != nil {
		if user.OwnsRelease(course.Release.ID) {
			owned = true
		}
	}

	metaTitle := strings.Trim(course.Title, " ")

	metaDescription := strings.Trim(course.Subtitle, " ")

	section := &db.Section{}
	if sectionIDParam != "" {
		var err2 error
		section, err2 = db.GetSection(sectionIDParam)
		if err2 == nil && section.Description != "" {
			metaDescription = "In this Section - " + strings.Trim(section.Description, " ")
			metaTitle += " - " + section.Name
		}
	}

	sectionMarkdownHTML := ""
	if section.GithubSection.MarkdownCache != "" {
		buffer, _ := markdown.Convert([]byte(section.GithubSection.MarkdownCache))
		sectionMarkdownHTML = buffer.String()
	}

	buffer, _ := markdown.Convert([]byte(course.Markdown))
	courseMarkdownHTML := buffer.String()

	c.HTML(
		http.StatusOK,
		"course_.html",
		gin.H{
			"Owned":               owned,
			"CourseMarkdownHTML":  template.HTML(courseMarkdownHTML),
			"SectionMarkdownHTML": template.HTML(sectionMarkdownHTML),
			"Section":             section,
			"User":                auth2.GetLoggedInUserLogError(c),
			"LoggedIn":            auth2.IsLoggedInValid(c),
			"Messages":            msg.GetMessages(c),
			"Releases":            releases,
			"Course":              course,
			"Meta": meta{
				Title:    metaTitle,
				Desc:     metaDescription,
				ImageURL: course.Release.ImageURL,
			},
		},
	)
}
