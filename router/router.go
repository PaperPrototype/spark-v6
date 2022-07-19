package router

import (
	"errors"
	"html/template"
	"log"
	"main/router/api"
	"main/router/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Setup() {
	router = gin.Default()

	// Logger middleware will write the logs to gin.DefaultWriter even when you set GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	// router.Use(gin.Logger())

	// recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		if c.Request.Host == "sparkv6.herokuapp.com" {
			log.Println("REDIRECTING to sparker3d.com")
			c.Redirect(http.StatusMovedPermanently, "https://sparker3d.com"+c.Request.URL.Path)
		}
	})

	router.RemoveExtraSlash = false
	router.RedirectTrailingSlash = true

	router.SetFuncMap(template.FuncMap{
		// a dictionary util that can be used to pass multiple inputs to a template
		// syntax:
		/*

			dict "key" .Value "key" .Value

		*/
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	})

	router.LoadHTMLGlob("./templates/*")
	router.Static("/resources", "./resources")

	router.StaticFile("favicon.ico", "./resources/images/favicon.ico")

	api.AddRoutes(router.Group("/api"))
	routes.AddRoutes(router)
}

func Run() {
	router.Run()
}
