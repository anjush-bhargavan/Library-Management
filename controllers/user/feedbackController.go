package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)


//Feedback handles users to give their feedback
func Feedback(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var feedback models.FeedBack

	if err :=c.ShouldBindJSON(&feedback);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"status":"Failed",
											"message":"Error while Binding",
											"data":err.Error(),
										})
		return
	}
	feedback.UserID=userID

	if err := config.DB.Create(&feedback).Error; err!= nil {
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Feedback sent succesfully",
						"data":feedback,
					})

}