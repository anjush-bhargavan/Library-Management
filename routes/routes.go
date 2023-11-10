package routes

import (
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)


//ConfigRoutes to handle routes
func ConfigRoutes(r *gin.Engine){

	r.Use(middleware.ClearCache())

	agentRoutes(r)
	userRoutes(r)
	adminRoutes(r)

}