package router2

import (
	"errors"
	"html/template"
	"log"
	v2 "main/router2/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
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

			{{ "my_template.html" dict "key" .Value "key" .Value }}

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

	router.LoadHTMLGlob("./templates2/*")
	router.Static("/resources2", "./resources2")
	router.StaticFile("favicon.ico", "./resources2/images/favicon.ico")

	// setup all the routes
	SetupRoutes()

	// v2 api for web app capabilities
	v2.AddRoutes(router.Group("/v2"))

	router.Run()
}
