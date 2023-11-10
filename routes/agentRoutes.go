package routes

import (
	agent "github.com/anjush-bhargavan/library-management/controllers/agent"
	"github.com/anjush-bhargavan/library-management/middleware"
	"github.com/gin-gonic/gin"
)


func agentRoutes(r *gin.Engine ) {

agentGroup := r.Group("/agent")
agentGroup.Use(middleware.Authorization("agent"))
{
	agentGroup.GET("/orders",agent.ViewOrders)
	agentGroup.GET("/order/:id",agent.GetOrder)
	agentGroup.POST("/order",agent.UpdateOrders)
}
}