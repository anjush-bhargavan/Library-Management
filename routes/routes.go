package routes

import (
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)



func RoutesConfig(r *gin.Engine){

	r.Use(middleware.ClearCache())

	userRoutes(r)
	adminRoutes(r)

}