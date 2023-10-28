package routes

import (
	"github.com/anjush-bhargavan/library-management/controllers"
	user "github.com/anjush-bhargavan/library-management/controllers/user"
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)


//function to handle user side routes
func userRoutes(r *gin.Engine){

	r.GET("/login",controllers.UserLoginPage)
	r.POST("/login",controllers.UserLogin)
	r.GET("/signup",controllers.UserSignupPage)
	r.POST("/signup",controllers.UserSignup)
	r.GET("/verifyotp",controllers.VerifyOTPPage)
	r.POST("/verifyotp",controllers.VerifyOTP)

	userGroup:=r.Group("/user")
	userGroup.Use(middleware.UserAuth())
	{
		userGroup.GET("/home",controllers.HomePage)
		userGroup.GET("home/book/:id",user.GetBook)
		userGroup.GET("home/viewbooks",user.ViewBooks)
		userGroup.GET("/logout",controllers.UserLogout)

		userGroup.GET("/profile")
		userGroup.GET("/profile/update",user.UserProfile)
		userGroup.PUT("/profile/update",user.ProfileUpdate)
		userGroup.GET("/profile/membership")
		userGroup.POST("/profile/memebership")
		userGroup.GET("/profile/viewfine")
		userGroup.POST("/profile/payfine")
		userGroup.GET("/profile/viewhistory")

		userGroup.POST("/addtocart/:id",user.AddToCart)
		userGroup.GET("/showcart",user.ShowCart)
		userGroup.DELETE("/deletecartitem/:id",user.DeleteCart)

	}
	

}	