package main

import (
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/routes"
	"github.com/gin-gonic/gin"
)



func main() {
	config.Loadenv()
	config.InitRedis()
	config.InitCron()
	r:=gin.Default()
	config.ConnectDB()
	r.LoadHTMLGlob("templates/*.html")
	routes.ConfigRoutes(r)

	r.Run("localhost:8080")
}