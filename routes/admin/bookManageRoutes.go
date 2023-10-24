package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//BookRoutes to handle book management on admin side
func BookRoutes(r *gin.Engine) {

	r.GET("/book/:id",admin.GetBook)
	r.POST("/addbooks",admin.AddBooks)
	r.GET("/viewbooks",admin.ViewBooks)
	r.PUT("/updatebook/:id",admin.UpdateBook)
	r.DELETE("/deletebook/:id",admin.DeleteBook)
}
