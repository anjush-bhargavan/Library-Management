package routes

import (
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
)

//CategoryRoutes to handle Category management on admin side
func CategoryRoutes(r *gin.Engine) {

	r.GET("/category/:id",admin.GetCategory)
	r.POST("/addcategorys",admin.AddCategorys)
	r.GET("/viewcategorys",admin.ViewCategorys)
	r.PUT("/updatecategory/:id",admin.UpdateCategory)
	r.DELETE("/deletecategory/:id",admin.DeleteCategory)
}