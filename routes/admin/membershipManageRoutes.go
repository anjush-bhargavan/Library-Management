package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//MembershipRoutes handles other routes of admin-side
func MembershipRoutes( r *gin.Engine) {
	r.GET("/memberships",admin.ViewMemberships)
	r.PATCH("/cancelmembership",admin.RemoveMembership)
	r.GET("/member/:id",admin.GetMembership)
}