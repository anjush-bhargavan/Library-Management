package routes

import (
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
)

//CategoryRoutes to handle Category management on admin side
func CategoryRoutes(r *gin.Engine) {

	r.GET("/category/:id",admin.GetCategory)
	r.POST("/category",admin.AddCategorys)
	r.GET("/categories",admin.ViewCategorys)
	r.PUT("/category/:id",admin.UpdateCategory)
	r.DELETE("/category/:id",admin.DeleteCategory)
}