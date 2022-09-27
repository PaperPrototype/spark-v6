package router2

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/db"
	"main/githubapi"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

func getMedia(c *gin.Context) {
	versionID := c.Params.ByName("versionID")
	mediaName := c.Params.ByName("mediaName")

	version, err1 := db.GetVersion(versionID)
	if err1 != nil {
		log.Println("routes/get ERROR getting version in getNameMedia:", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	/* TODO?
	maybe not necessary since course versions can't be viewed unless the course is free (or the user has paid)
	and getting access to the image links without access to the course would be difficult
	*/
	// check if course release is free
	// if paid
	//	 check if student has access to course
	// else
	// 	 free so anyone can view it?

	// if it is a github based version
	if version.HasGithubVersion() {
		githubVersion, err2 := version.GetGithubVersion()
		if err2 != nil {
			log.Println("routes/get ERROR getting githubVersion in getNameMedia:", err2)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err3 := version.GetAuthorUser()
		if err3 != nil {
			log.Println("routes/get ERROR getting githubVersion in getNameMedia:", err2)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		githubConnection, err4 := githubapi.GetGithubConnection(user)
		if err4 != nil {
			log.Println("routes/get ERROR getting authors github connection in getNameMedia", err4)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		client := githubapi.NewClient(githubConnection, ctx)

		githubUser, _, err5 := client.Users.Get(ctx, "")
		if err5 != nil {
			log.Println("routes/get ERROR getting githubUser in getNameMedia", err5)
			c.AbortWithStatus(http.StatusNotFound) // user should not know of the existence of this file
			return
		}

		repo, _, err6 := client.Repositories.GetByID(ctx, githubVersion.RepoID)
		if err6 != nil {
			log.Println("routes/get ERROR getting repo by ID in getNameMedia", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		mediaType := filepath.Ext(mediaName)

		if mediaType == ".zip" {
			readCloser, err7 := client.Repositories.DownloadContents(ctx, *githubUser.Login, *repo.Name, "Resources/"+mediaName, &github.RepositoryContentGetOptions{
				Ref: githubVersion.SHA,
			})
			if err7 != nil {
				log.Println("routes/get ERROR getting downloading contents in getNameMedia", err6)
				return
			}
			defer readCloser.Close()

			written, err8 := io.Copy(c.Writer, readCloser)
			if err8 != nil {
				log.Println("routes/get ERROR copying/writing contents in getNameMedia", err6)
				return
			}

			c.Writer.Header().Set("Content-Type", mediaType)
			c.Writer.Header().Set("Content-Length", fmt.Sprint(written))
			return
		}

		readCloser, err7 := client.Repositories.DownloadContents(ctx, *githubUser.Login, *repo.Name, "Assets/"+mediaName, &github.RepositoryContentGetOptions{
			Ref: githubVersion.SHA,
		})
		if err7 != nil {
			log.Println("routes/get ERROR getting downloading contents in getNameMedia", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer readCloser.Close()

		written, err8 := io.Copy(c.Writer, readCloser)
		if err8 != nil {
			log.Println("routes/get ERROR copying/writing contents in getNameMedia", err6)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Writer.Header().Set("Content-Type", mediaType)
		c.Writer.Header().Set("Content-Length", fmt.Sprint(written))
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
