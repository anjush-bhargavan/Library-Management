package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)


func UserManageRoutes(c *gin.Engine) {
	c.GET("/user/:id",admin.ViewUser)
	c.POST("/adduser",admin.AddUser)
	c.PUT("/updateuser/:id",admin.UpdateUser)
	c.DELETE("/deleteuser/:id",admin.DeleteUser)

}