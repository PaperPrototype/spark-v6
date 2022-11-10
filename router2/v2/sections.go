package v2

import (
	"fmt"
	"log"
	"main/auth2"
	"main/db"
	"main/githubapi"
	"main/markdown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// create or update github section
func postSectionGithubFORM(c *gin.Context) {
	sectionID := c.Params.ByName("sectionID")
	path := c.PostForm("path")

	if path == "" {
		c.JSON(http.StatusBadRequest, payload{
			Error: "No markdown file selected. Path cannot be empty.",
		})
		return
	}

	sectionIDNum, _ := strconv.ParseUint(sectionID, 10, 64)

	section, err := db.GetSection(sectionID)
	if err != nil {
		log.Println("v2/section.go ERROR getting section in postSectionGithubFORM:", err)
		c.JSON(http.StatusBadRequest, payload{
			Error: "Error getting section.",
		})
		return
	}

	// get markdown so we can cache it
	author := auth2.GetLoggedInUserLogError(c)
	release, _ := db.GetAnyRelease(section.ReleaseID)
	markdown, _ := githubapi.GetGithubMarkdown(author, release, path)

	// if no github section
	if section.ID != section.GithubSection.SectionID {
		// create github section

		githubSection := db.GithubSection{
			SectionID:     sectionIDNum,
			Path:          path,
			MarkdownCache: markdown,
		}
		err := db.CreateGithubSection(sectionID, &githubSection)
		if err != nil {
			log.Println("v2/section.go ERROR creating github section in postSectionGithubFORM:", err)
			c.JSON(http.StatusBadRequest, payload{
				Error: "Error creating github section.",
			})
			return
		}
	} else {
		// update github section

		err := db.UpdateGithubSection(sectionID, path, markdown)
		if err != nil {
			log.Println("v2/section.go ERROR updating github section in postSectionGithubFORM:", err)
			c.JSON(http.StatusBadRequest, payload{
				Error: "Error updating github section.",
			})
			return
		}
	}

	c.JSON(http.StatusOK, payload{})
}

func getSectionJSON(c *gin.Context) {
	sectionID := c.Params.ByName("sectionID")

	section, err := db.GetSection(sectionID)
	if err != nil {
		log.Println("v2/sections.go ERROR getting section in getSection:", err)
		c.JSON(http.StatusOK, payload{
			Error:   "Error getting section",
			Payload: section,
		})
		return
	}

	c.JSON(http.StatusOK, payload{
		Payload: section,
	})
}

func getSectionMarkdownHTML(c *gin.Context) {
	sectionID := c.Params.ByName("sectionID")

	section, err := db.GetSection(sectionID)
	if err != nil {
		log.Println("v2/sections.go ERROR getting section in getSection:", err)
		c.JSON(http.StatusOK, payload{
			Error:   "Error getting section",
			Payload: section,
		})
		return
	}

	release, err1 := db.GetAnyRelease(section.ReleaseID)
	if err1 != nil {
		log.Println("v2/middlewares.go ERROR getting any release in mustBeAuthorSectionID:", err1)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting release for that section",
		})
		return
	}

	course, err2 := db.GetCourse(release.CourseID)
	if err2 != nil {
		log.Println("v2/middlewares.go ERROR getting course in mustBeAuthorSectionID:", err2)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting course for that section",
		})
		return
	}

	author, err10 := db.GetUser(course.UserID)
	if err10 != nil {
		log.Println("v2/middlewares.go ERROR getting user in mustBeAuthorSectionID:", err2)
		c.JSON(http.StatusOK, payload{
			Error: "Error getting course author",
		})
		return
	}

	loggedIn := false
	if auth2.IsLoggedInValid(c) {
		loggedIn = true
	}

	user := &db.User{}
	if loggedIn {
		user = auth2.GetLoggedInUserLogError(c)
	}

	// if release is private
	if !release.Public {
		if !loggedIn {
			c.JSON(http.StatusNotFound, payload{})
			return
		} else if course.UserID != user.ID {
			c.JSON(http.StatusNotFound, payload{})
			return
		}
	}

	// if not the author
	if course.UserID != user.ID {
		// if section is not set as "free preview section"
		if !section.Free {
			// if paid course release
			if !release.IsFree() {
				// if does not own the course release
				if !user.OwnsRelease(release.ID) {
					if !loggedIn {
						// return login message
						c.JSON(http.StatusOK, payload{
							Payload: `
							<div x-data style="text-align:center; margin-top:40vh;">
								<i class="fa-solid fa-lock" style="font-size:2rem; margin:2rem;"></i>
								<p style="padding:2rem;">Login or signup to unlock all sections of the course</p>
								<button @click="onboard_open()" class="thm-bg-hl utls-bd" style="padding:1rem;">Login or Signup</button>
							</div>`,
						})
						return
					}

					// return buy course message
					c.JSON(http.StatusOK, payload{
						Payload: `
					<div style="text-align:center; margin-top:40vh;">
						<i class="fa-solid fa-lock" style="font-size:2rem; margin:2rem;"></i>
						<p style="padding:2rem;">Buy ` + course.Title + ` to unlock all sections of the course. One time payment</p>
						<a href="/` + author.Username + `/` + course.Name + `/buy/` + fmt.Sprint(release.ID) + `">
							<button class="thm-bg-hl utls-bd" style="padding:1rem;">USD $` + fmt.Sprint(release.GetPriceUSD()) + `.00</button>
						</a>
					</div>`,
					})
					return
				}

				// section is not free, and course is free
				// and not logged in
				// must log in to view section
			} else if !loggedIn {
				// return login message
				c.JSON(http.StatusOK, payload{
					Payload: `
					<div x-data style="text-align:center; margin-top:40vh;">
						<i class="fa-solid fa-lock" style="font-size:2rem; margin:2rem;"></i>
						<p style="padding:2rem;">Login or signup to unlock all sections of this course for free</p>
						<button @click="onboard_open()" class="thm-bg-hl utls-bd" style="padding:1rem;">Login or Signup</button>
					</div>`,
				})
				return
			}
		}
	}

	html := `<p style="text-align:center;">This section is empty</p>`

	// if author
	if user.ID == course.UserID {
		html = `<p style="text-align:center;">Open course settings to connect to a github release.</p>`
	}

	// for github based courses
	if release.HasGithubRelease() {
		// has github section
		if section.GithubSection.SectionID == section.ID {
			tmpMarkdown := ""
			Error := ""

			// if markdown IS cached
			if section.GithubSection.MarkdownCache != "" {
				// if cache is invalid (github release was updated)
				if section.GithubSection.MarkdownCachePatchNum != release.GithubRelease.Patch {
					// cache the markdown from this section if not cached already
					tmpMarkdown, Error = githubapi.GetGithubMarkdown(author, release, section.GithubSection.Path)
					db.UpdateGithubSectionMarkdownCache(section.ID, tmpMarkdown)
				} else {
					// ALREADY CACHED!
					tmpMarkdown = section.GithubSection.MarkdownCache
				}

			} else {
				// cache the markdown from this section if not cached already
				tmpMarkdown, Error = githubapi.GetGithubMarkdown(author, release, section.GithubSection.Path)
				db.UpdateGithubSectionMarkdownCache(section.ID, tmpMarkdown)
			}

			// if tmpMarkdown not empty
			if tmpMarkdown != "" {
				buffer, err9 := markdown.Convert([]byte(tmpMarkdown))
				if err9 != nil {
					log.Println("api/sections ERROR converting content tp markdown in getSectionMarkdownHTML:", err9)
				}

				// set HTML variable
				html = buffer.String()
			}

			c.JSON(http.StatusOK, payload{
				Error:   Error,
				Payload: html,
			})
			return
		}

		// if author
		if user.ID == course.UserID {
			html = `<p style="text-align:center;">Open this sections settings to select a markdown file from the github repository.</p>`
		}
	}

	// else TODO get html contents from user created (non github based) course

	// return html
	c.JSON(http.StatusOK, payload{
		Payload: html,
	})
}

func postSection(c *gin.Context) {
	id := c.Params.ByName("sectionID")
	name := c.PostForm("name")
	desc := c.PostForm("desc")
	free := c.PostForm("free")
	num := c.PostForm("num")

	log.Println("free was json as:", free)

	isFree := false
	if free == "true" {
		isFree = true
	}

	numInt16, _ := strconv.ParseUint(num, 10, 64)
	err := db.UpdateSection(id, name, desc, isFree, uint16(numInt16))
	if err != nil {
		c.JSON(http.StatusOK, payload{
			Error: "Error updating section",
		})
		return
	}

	c.JSON(http.StatusOK, payload{})
}

func deleteSection(c *gin.Context) {
	id := c.Params.ByName("sectionID")

	err := db.DeleteSection(id)
	if err != nil {
		log.Println("v2/sections.go ERROR deleting section in deleteSection:", err)
		c.JSON(http.StatusOK, payload{
			Error: "Error deleting section",
		})
		return
	}

	c.JSON(http.StatusOK, payload{})
}
