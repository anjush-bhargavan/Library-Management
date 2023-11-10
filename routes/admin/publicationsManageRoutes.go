package routes


import (
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
)

//PublicationRoutes to handle Publication management on admin side
func PublicationRoutes(r *gin.Engine) {

	r.GET("/publication/:id",admin.GetPublication)
	r.POST("/publication",admin.AddPublication)
	r.GET("/publications",admin.ViewPublications)
	r.PUT("/publication/:id",admin.UpdatePublication)
	r.DELETE("/publication/:id",admin.DeletePublication)
}