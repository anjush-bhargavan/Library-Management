package routes

import (
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)



func AdminRoutes(r *gin.Engine) {

	r.Use(middleware.AdminAuth())

	bookRoutes(r)
}