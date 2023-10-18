package main

import (
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/routes"
	"github.com/gin-gonic/gin"
)



func main() {
	r:=gin.Default()
	config.ConnectDB()

	routes.RoutesConfig(r)

	r.Run("localhost:8080")
}