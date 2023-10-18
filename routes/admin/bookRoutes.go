package routes

import (
	"github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)


func bookRoutes(r *gin.Engine) {

	r.POST("/addbooks",controllers.AddBooks)
	r.GET("/viewbooks",controllers.ViewBooks)
	r.PUT("/book/:id",controllers.UpdateBook)
	r.DELETE("/book/:id",controllers.DeleteBook)
}
