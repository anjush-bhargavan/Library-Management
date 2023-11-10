package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//UserManageRoutes to handle user management on admin side
func UserManageRoutes(c *gin.Engine) {
	c.GET("/user/:id",admin.ViewUser)
	c.GET("/users",admin.ViewUsers)
	c.PUT("/user/:id",admin.UpdateUser)
	c.DELETE("/user/:id",admin.DeleteUser)

}