package routes

import (
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
	"github.com/gin-gonic/gin"
)

//RecordRoutes handles other routes of admin-side
func RecordRoutes( r *gin.Engine) {
	r.GET("/history",admin.ViewHistory)
	r.GET("/booksout",admin.ViewBooksOut)
	r.GET("/orders",admin.ViewOrders)
}