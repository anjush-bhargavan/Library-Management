package routes


import (
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
)

//PublicationRoutes to handle Publication management on admin side
func PublicationRoutes(r *gin.Engine) {

	r.GET("/publication/:id",admin.GetPublication)
	r.POST("/addpublications",admin.AddPublication)
	r.GET("/viewpublications",admin.ViewPublications)
	r.PUT("/updatepublication/:id",admin.UpdatePublication)
	r.DELETE("/deletepublication/:id",admin.DeletePublication)
}