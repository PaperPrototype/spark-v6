package v2

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/auth2"
	"main/db"
	"main/githubapi"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

// allow anyone to view images from the course
func getReleaseGithubAsset(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")
	name := c.Params.ByName("name")

	release, err := db.GetAnyRelease(releaseID)
	if err != nil {
		log.Println("v2/releases.go ERROR getting release in getReleaseGithubAsset:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if release.GithubEnabled {
		githubRelease, err1 := db.GetGithubReleaseWithIDStr(releaseID)
		if err1 != nil {
			log.Println("v2/releases.go ERROR getting github release in getReleaseGithubAsset:", err1)
			c.Status(http.StatusInternalServerError)
			return
		}

		user, err3 := release.GetAuthorUser()
		if err3 != nil {
			log.Println("v2/releases.go ERROR getting githubVersion in getReleaseGithubAsset:", err3)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		githubConnection, err4 := githubapi.GetGithubConnection(user)
		if err4 != nil {
			log.Println("v2/releases.go ERROR getting authors github connection in getReleaseGithubAsset", err4)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		client := githubapi.NewClient(githubConnection, ctx)

		githubUser, _, err5 := client.Users.Get(ctx, "")
		if err5 != nil {
			log.Println("v2/releases.go ERROR getting githubUser in getReleaseGithubAsset", err5)
			c.AbortWithStatus(http.StatusNotFound) // user should not know of the existence of this file
			return
		}

		repo, _, err6 := client.Repositories.GetByID(ctx, int64(githubRelease.RepoID))
		if err6 != nil {
			log.Println("v2/releases.go ERROR getting repo by ID in getReleaseGithubAsset", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		mediaType := filepath.Ext(name)

		if mediaType != ".md" {
			readCloser, err7 := client.Repositories.DownloadContents(ctx, *githubUser.Login, *repo.Name, "Assets/"+name, &github.RepositoryContentGetOptions{
				Ref: githubRelease.SHA,
			})
			if err7 != nil {
				log.Println("v2/releases.go ERROR getting downloading contents in getReleaseGithubAsset", err6)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer readCloser.Close()

			written, err8 := io.Copy(c.Writer, readCloser)
			if err8 != nil {
				log.Println("v2/releases.go ERROR copying/writing contents in getReleaseGithubAsset", err6)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Writer.Header().Set("Content-Type", mediaType)
			c.Writer.Header().Set("Content-Length", fmt.Sprint(written))
			return
		}

		c.Status(http.StatusNoContent)
	} else {
		// TODO user uploaded courses and non github courses
	}

	c.Status(http.StatusNoContent)
}

func getGithubReleaseAssetsTreeJSON(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	githubRelease, err := db.GetGithubReleaseWithIDStr(releaseID)
	if err != nil {
		log.Println("v2/github.go ERROR getting github release in getGithubReleaseTreeJSON:", err)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github release",
			},
		)
		return
	}

	user := auth2.GetLoggedInUserLogError(c)

	connection, err1 := githubapi.GetGithubConnection(user)
	if err1 != nil {
		log.Println("v2/github.go ERROR getting github connection in getGithubReleaseTreeJSON:", err1)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github connection",
			},
		)
		return
	}

	ctx := context.Background()

	// get client
	client := githubapi.NewClient(connection, ctx)

	githubUser, _, err2 := client.Users.Get(ctx, "")
	if err2 != nil {
		log.Println("api/github ERROR getting github user in getGithubReleaseTreeJSON:", err2)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github user",
			},
		)
		return
	}

	repo, _, err3 := client.Repositories.GetByID(ctx, int64(githubRelease.RepoID))
	if err3 != nil {
		log.Println("api/github ERROR getting repo by ID in getGithubReleaseTreeJSON:", err3)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting repository",
			},
		)
		return
	}

	// get folders from repo with info from githubVersion
	// use sha to get specific commit
	tree, _, err4 := client.Git.GetTree(ctx, *githubUser.Login, *repo.Name, githubRelease.SHA, true)
	if err4 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubReleaseTreeJSON:", err4)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting repository tree",
			},
		)
		return
	}

	entrees := []github.TreeEntry{}
	for _, entree := range tree.Entries {
		if strings.Contains(*entree.Path, ".jpg") || strings.Contains(*entree.Path, ".png") ||
			strings.Contains(*entree.Path, ".jpeg") || strings.Contains(*entree.Path, ".gif") {
			entrees = append(entrees, entree)
		}
	}

	c.JSON(
		http.StatusOK,
		payload{
			Payload: entrees,
		},
	)
}

// get the files tree of the github repository connected to the release
func getGithubReleaseTreeJSON(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	githubRelease, err := db.GetGithubReleaseWithIDStr(releaseID)
	if err != nil {
		log.Println("v2/github.go ERROR getting github release in getGithubReleaseTreeJSON:", err)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github release",
			},
		)
		return
	}

	user := auth2.GetLoggedInUserLogError(c)

	connection, err1 := githubapi.GetGithubConnection(user)
	if err1 != nil {
		log.Println("v2/github.go ERROR getting github connection in getGithubReleaseTreeJSON:", err1)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github connection",
			},
		)
		return
	}

	ctx := context.Background()

	// get client
	client := githubapi.NewClient(connection, ctx)

	githubUser, _, err2 := client.Users.Get(ctx, "")
	if err2 != nil {
		log.Println("api/github ERROR getting github user in getGithubReleaseTreeJSON:", err2)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting github user",
			},
		)
		return
	}

	repo, _, err3 := client.Repositories.GetByID(ctx, int64(githubRelease.RepoID))
	if err3 != nil {
		log.Println("api/github ERROR getting repo by ID in getGithubReleaseTreeJSON:", err3)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting repository",
			},
		)
		return
	}

	// get folders from repo with info from githubVersion
	// use sha to get specific commit
	tree, _, err4 := client.Git.GetTree(ctx, *githubUser.Login, *repo.Name, githubRelease.SHA, true)
	if err4 != nil {
		log.Println("api/github ERROR getting repo contents in getGithubReleaseTreeJSON:", err4)
		c.JSON(
			http.StatusInternalServerError,
			payload{
				Error: "Error getting repository tree",
			},
		)
		return
	}

	entrees := []github.TreeEntry{}
	for _, entree := range tree.Entries {
		if strings.Contains(*entree.Path, ".md") {
			entrees = append(entrees, entree)
		}
	}

	c.JSON(
		http.StatusOK,
		payload{
			Payload: entrees,
		},
	)
}

func postReleaseFORM(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	priceStr := c.PostForm("price")
	publicStr := c.PostForm("public")
	githubEnabledStr := c.PostForm("githubEnabled")
	postsNeededStr := c.PostForm("postsNeededNum")

	price, _ := strconv.ParseUint(priceStr, 10, 64)
	postsNeededNum, _ := strconv.ParseUint(postsNeededStr, 10, 64)
	public := publicStr == "true"
	githubEnabled := githubEnabledStr == "true"
	imageURL := c.PostForm("imageURL")

	fixedImageURL := ""
	if len(imageURL) > 0 {
		fixedImageURL = "/v2/releases/" + releaseID + "/assets/" + imageURL[7:len(imageURL)]
	}

	db.UpdateRelease(releaseID, price*100, public, uint16(postsNeededNum), fixedImageURL, githubEnabled)

	c.JSON(http.StatusOK, payload{})
}

func getReleaseJSON(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	message := ""

	release, err := db.GetAnyRelease(releaseID)
	if err != nil {
		message = "Error getting release"
	}

	c.JSON(http.StatusOK, payload{
		Error:   message,
		Payload: *release,
	})
}

// get a github release
func getGithubReleaseJSON(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	message := ""

	githubRelease, err := db.GetGithubReleaseWithIDStr(releaseID)
	if err != nil {
		// "please select a github repository since this endpoint is used when selecting a repository in case if one doesn't exist yet"
		message = "No github release. Please select a github repository."
	}

	c.JSON(http.StatusOK, payload{
		Error:   message,
		Payload: *githubRelease,
	})
}

// update or create a github release
func postGithubReleaseFORM(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	repoID := c.PostForm("repoID")
	branch := c.PostForm("branch")
	SHA := c.PostForm("SHA")

	releaseIDNum, _ := strconv.ParseUint(releaseID, 10, 64)
	repoIDNum, _ := strconv.ParseUint(repoID, 10, 64)

	githubRelease := db.GithubRelease{
		ReleaseID: releaseIDNum,
		RepoID:    repoIDNum,
		Branch:    branch,
		SHA:       SHA,
	}

	// releaseID is used to check if a githubRelease already exists or not
	err := db.CreateOrUpdateGithubRelease(releaseID, &githubRelease)
	if err != nil {
		log.Println("v2/release.go ERROR upating or creating github release:", err)
		c.JSON(http.StatusInternalServerError, payload{
			Error: "An error occured updating or creating github release",
		})
		return
	}

	c.JSON(http.StatusOK, payload{})
}

// get sections for a release
func getReleaseSectionsJSON(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	sections, err := db.GetReleaseSections(releaseID)
	if err != nil {
		log.Println("v2/release.go ERROR getting release sections:", err)
		c.JSON(http.StatusInternalServerError, payload{
			Error: "Error getting release sections",
		})
		return
	}

	c.JSON(http.StatusOK, payload{
		Payload: sections,
	})
}

// create a new section
func postReleaseSectionFORM(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	releaseIDNum, err := strconv.ParseUint(releaseID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, payload{
			Error: "Parse error for releaseID",
		})
		return
	}

	num := uint16(0)
	sections, err2 := db.GetReleaseSections(releaseID)
	if err2 == nil && len(sections) > 0 {
		// get the greatest number currently and set ours to it plus 1
		num = sections[len(sections)-1].Num + 1
	}

	name := c.PostForm("name")
	section := db.Section{
		Num:       num,
		Name:      name,
		ReleaseID: releaseIDNum,
	}
	err1 := db.CreateSection(&section)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, payload{
			Error: "Error creating section",
		})
		return
	}

	// must always respond with json
	c.JSON(http.StatusOK, payload{})
}

/*



























 */

func getSections(tree github.Tree) []db.Section {
	sections := []db.Section{}
	for _, entry := range tree.Entries {
		// no sub sections
		if sectionNameAllowed(*entry.Path) {
			sections = append(sections, db.Section{
				Name: *entry.Path,
			})
		}
	}

	return sections
}

func sectionNameAllowed(str string) bool {
	lowerStr := strings.ToLower(str)
	if strings.Contains(lowerStr, "/") || strings.Contains(lowerStr, "ignore") || strings.Contains(lowerStr, "resources") || strings.Contains(lowerStr, "assets") {
		return false
	}

	return true
}

func deleteRelease(c *gin.Context) {
	releaseID := c.Params.ByName("releaseID")

	err := db.DeleteRelease(releaseID)
	if err != nil {
		log.Println("v2/releases.go ERROR deleting release in deleteRelease:", err)
		c.JSON(http.StatusOK, payload{
			Error: "Error deleting release",
		})
		return
	}

	c.JSON(http.StatusOK, payload{})
}
