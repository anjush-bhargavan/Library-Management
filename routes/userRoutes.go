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
		userGroup.GET("/logout",)
	}
	

}	