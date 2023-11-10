package routes

import (
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
)

//AuthorRoutes to handle Author management on admin side
func AuthorRoutes(r *gin.Engine) {

	r.GET("/author/:id",admin.GetAuthor)
	r.POST("/author",admin.AddAuthors)
	r.GET("/authors",admin.ViewAuthors)
	r.PUT("/author/:id",admin.UpdateAuthor)
	r.DELETE("/author/:id",admin.DeleteAuthor)
}