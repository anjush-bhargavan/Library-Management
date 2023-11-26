package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//OtherRoutes handles other routes of admin-side
func OtherRoutes( r *gin.Engine) {
	r.GET("/feedbacks",admin.ViewFeedbacks)
}