package routes

import "github.com/gin-gonic/gin"


func profileRoutes(c *gin.Engine) {

	c.GET("/profile")
	c.GET("/profile/update")
	c.PUT("/profile/update")
	c.GET("/profile/membership")
	c.POST("/profile/memebership")
	c.GET("/profile/viewfine")
	c.POST("/profile/payfine")
	c.GET("/profile/viewhistory")
}