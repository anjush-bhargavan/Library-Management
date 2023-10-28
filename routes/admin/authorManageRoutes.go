package routes

import (
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/controllers/admin"
)

//AuthorRoutes to handle Author management on admin side
func AuthorRoutes(r *gin.Engine) {

	r.GET("/author/:id",admin.GetAuthor)
	r.POST("/addauthors",admin.AddAuthors)
	r.GET("/viewauthors",admin.ViewAuthors)
	r.PUT("/updateauthor/:id",admin.UpdateAuthor)
	r.DELETE("/deleteauthor/:id",admin.DeleteAuthor)
}