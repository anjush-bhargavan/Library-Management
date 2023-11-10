package main

import (
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/routes"
	"github.com/gin-gonic/gin"
)



func main() {
	config.Loadenv()
	r:=gin.Default()
	config.ConnectDB()
	config.InitRedis()
	r.LoadHTMLGlob("templates/*.html")
	routes.ConfigRoutes(r)

	r.Run("localhost:8080")
}