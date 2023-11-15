package routes

import (
	"github.com/anjush-bhargavan/library-management/controllers"
	user "github.com/anjush-bhargavan/library-management/controllers/user"
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)

// function to handle user side routes
func userRoutes(r *gin.Engine) {

	r.GET("/login", controllers.UserLoginPage)
	r.POST("/login", controllers.UserLogin)
	r.GET("/signup", controllers.UserSignupPage)
	r.POST("/signup", controllers.UserSignup)
	r.GET("/verifyotp", controllers.VerifyOTPPage)
	r.POST("/verifyotp", controllers.VerifyOTP)

	userGroup := r.Group("/user")
	userGroup.Use(middleware.Authorization("user"))
	{
		userGroup.GET("/home", controllers.HomePage)
		userGroup.GET("/logout", controllers.UserLogout)

		userGroup.GET("home/book/:id", user.GetBook)
		userGroup.GET("home/books", user.ViewBooks)
		userGroup.GET("/search",user.SearchBooks)

		userGroup.GET("/profile", user.UserProfile)
		userGroup.PUT("/profile", user.ProfileUpdate)
		userGroup.PATCH("/profile/changepassword", user.ChangePassword)

		userGroup.GET("profile/plans", user.ShowPlans)
		userGroup.POST("profile/plans", user.GetPlan)

		userGroup.GET("/profile/viewfine",user.ViewFine)
		r.GET("/profile/payfine",user.PayFine)
		
		userGroup.GET("/profile/viewhistory", user.ViewHistory)

		userGroup.POST("/wishlist/:id", user.AddToWishlist)
		userGroup.GET("/wishlist", user.ShowWishlist)
		userGroup.DELETE("/wishlist/:id", user.DeleteWishlist)

		r.GET("/profile/membership", user.Membership)
		r.GET("/payment/success", user.RazorpaySuccess)
		r.GET("/success", user.SuccessPage)
		r.GET("/invoice/download")

		userGroup.GET("/checkout/:id", user.DeliveryDetails)
		userGroup.POST("/checkout", user.Delivery)
		userGroup.PATCH("/cancel",user.CancelOrder)

		userGroup.POST("/review/:id",user.AddReview)
		userGroup.GET("review/:id",user.ShowReview)
		userGroup.PATCH("review/:id",user.EditReview)
		userGroup.DELETE("review/:id",user.DeleteReview)

		userGroup.POST("/feedback",user.Feedback)
		userGroup.GET("/events",user.ViewEvents)

	}

}
