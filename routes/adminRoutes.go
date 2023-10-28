package routes

import (
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
	admin "github.com/anjush-bhargavan/library-management/routes/admin"
)


//function to handle admin side routes
func adminRoutes(r *gin.Engine) {

	r.Use(middleware.AdminAuth())

	admin.BookRoutes(r)
	admin.UserManageRoutes(r)
	admin.AuthorRoutes(r)
	admin.CategoryRoutes(r)
	admin.PublicationRoutes(r)
}