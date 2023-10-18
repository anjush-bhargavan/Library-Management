package routes

import (
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/routes/admin"
)



func RoutesConfig(r *gin.Engine){

	r.Use(middleware.ClearCache())

	userRoutes(r)
	admin.AdminRoutes(r)

}