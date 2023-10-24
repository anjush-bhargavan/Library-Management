package routes

import (
	"github.com/anjush-bhargavan/library-management/controllers"
	user "github.com/anjush-bhargavan/library-management/controllers/user"
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)



func userRoutes(r *gin.Engine){

	r.GET("/login",controllers.UserLoginPage)
	r.POST("/login",controllers.UserLogin)
	r.GET("/signup",controllers.UserSignupPage)
	r.POST("/signup",controllers.UserSignup)

	userGroup:=r.Group("/user")
	userGroup.Use(middleware.UserAuth())
	{
		userGroup.GET("/home",controllers.HomePage)
		userGroup.GET("home/book/:id",user.GetBook)
		userGroup.GET("home/viewbooks",user.ViewBooks)
		userGroup.GET("/logout",controllers.UserLogout)

		userGroup.GET("/profile")
		userGroup.GET("/profile/update")
		userGroup.PUT("/profile/update")
		userGroup.GET("/profile/membership")
		userGroup.POST("/profile/memebership")
		userGroup.GET("/profile/viewfine")
		userGroup.POST("/profile/payfine")
		userGroup.GET("/profile/viewhistory")
	}
	

}	