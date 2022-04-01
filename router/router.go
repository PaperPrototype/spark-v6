package router

import (
	"errors"
	"html/template"
	"main/router/api"
	"main/router/routes"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Setup() {
	router = gin.Default()
	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = true

	router.SetFuncMap(template.FuncMap{
		// a sictionary util that can be used to pass input to templates
		// much like gin.H{}
		// gotten from
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

	api.AddRoutes(router.Group("/api"))
	routes.AddRoutes(router)
}

func Run() {
	router.Run()
}
