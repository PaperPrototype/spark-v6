package router

import (
	"main/router/api"
	"main/router/routes"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Setup() {
	router = gin.Default()

	router.LoadHTMLGlob("./templates/*")
	router.Static("/resources", "./resources")

	api.AddRoutes(router.Group("/api"))
	routes.AddRoutes(router)
}

func Run() {
	router.Run()
}
