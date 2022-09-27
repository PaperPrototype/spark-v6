package v2

import (
	"context"
	"fmt"
	"log"
	"main/auth2"
	"main/db"
	"main/githubapi"
	"main/markdown"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
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

	githubSection := db.GithubSection{
		SectionID: sectionIDNum,
		Path:      path,
	}
	err := db.CreateOrUpdateGithubSection(sectionID, &githubSection)
	if err != nil {
		log.Println("v2/section.go ERROR creating or updating section in postSectionGithubFORM:", err)
		c.JSON(http.StatusBadRequest, payload{
			Error: "Error creating or updating github section.",
		})
		return
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
		// if paid course release
		if !release.IsFree() {
			// if section is not set as "free preview section"
			if !section.Free {
				// if does not own the course release
				if !user.OwnsRelease(release.ID) {
					if !loggedIn {
						// return login to buy message
						c.JSON(http.StatusOK, payload{
							Payload: `
						<div x-data style="text-align:center; margin-top:40vh;">
							<p>Login or signup to buy ` + course.Title + ` and unlock all sections of the course.</p>
							<button @click="onboard_open()" class="thm-bg-hl utls-bd" style="padding:1rem;">Login or Signup</button>
						</div>`,
						})
						return
					}

					// return buy course message
					c.JSON(http.StatusOK, payload{
						Payload: `
					<div style="text-align:center; margin-top:40vh;">
						<p>Buy ` + course.Title + ` to unlock all sections of the course. One time payment.</p>
						<a href="/` + author.Username + `/` + course.Name + `/buy/` + fmt.Sprint(release.ID) + `">
							<button class="thm-bg-hl utls-bd" style="padding:1rem;">USD $` + fmt.Sprint(release.GetPriceUSD()) + `.00</button>
						</a>
					</div>`,
					})
					return
				}
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
			// get authors's github connection
			connection, err4 := githubapi.GetGithubConnection(author)
			if err4 != nil {
				log.Println("v2/sections.go ERROR getting user's github connection in getSectionMarkdown:", err4)

				c.JSON(http.StatusOK, payload{
					Error:   "Error getting github connection for author",
					Payload: html,
				})
				return
			}

			ctx := context.Background()

			// get client
			client := githubapi.NewClient(connection, ctx)

			githubUser, _, err5 := client.Users.Get(ctx, "")
			if err5 != nil {
				log.Println("v2/sections.go ERROR getting github user in getSectionMarkdown:", err5)

				c.JSON(http.StatusOK, payload{
					Error:   "Error getting github user",
					Payload: html,
				})
				return
			}

			repo, _, err6 := client.Repositories.GetByID(ctx, int64(release.GithubRelease.RepoID))
			if err6 != nil {
				log.Println("v2/sections.go ERROR getting repo by ID in getSectionMarkdown:", err6)

				c.JSON(http.StatusOK, payload{
					Error:   "Error getting github user",
					Payload: html,
				})
				return
			}

			// get folders from repo with info from githubVersion
			// use sha to get specific commit
			contentEncoded, _, _, err7 := client.Repositories.GetContents(ctx, *githubUser.Login, *repo.Name, section.GithubSection.Path, &github.RepositoryContentGetOptions{
				Ref: release.GithubRelease.SHA,
			})
			if err7 != nil {
				log.Println("api/github ERROR getting repo contents in getGithubRepoCommitContent:", err7)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// decode content
			content, err8 := contentEncoded.GetContent()
			if err8 != nil {
				log.Println("api/github ERROR decoding", *contentEncoded.Encoding, "content in getGithubRepoCommitContent:", err8)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			buffer, err9 := markdown.Convert([]byte(content))
			if err9 != nil {
				log.Println("api/github ERROR converting content ot markdown in getGithubRepoCommitContent:", err9)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.JSON(http.StatusOK, payload{
				Payload: buffer.String(),
			})
			return
		}

		// if author
		if user.ID == course.UserID {
			html = `<p style="text-align:center;">Open this sections settings to select a markdown file from the github repository.</p>`
		}
	}

	// else TODO get html contents from user created (non github based) course
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
