package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//EventRoutes handles other routes of admin-side
func EventRoutes( r *gin.Engine) {
	r.POST("/event",admin.AddEvent)
	r.PATCH("/event",admin.EditEvent)
	r.GET("/events",admin.ViewEvents)
}