package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//BookRoutes to handle book management on admin side
func BookRoutes(r *gin.Engine) {

	r.GET("/book/:id",admin.GetBook)
	r.POST("/book",admin.AddBooks)
	r.GET("/books",admin.ViewBooks)
	r.PUT("/book/:id",admin.UpdateBook)
	r.DELETE("/book/:id",admin.DeleteBook)
	r.DELETE("/review",admin.DeleteReview)
}
