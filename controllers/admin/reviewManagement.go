package controllers

import (
	"net/http"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//DeleteReview handles admin to delete the review of users using query by url
func DeleteReview(c *gin.Context) {
	 bookID:= c.Query("book_id")
	 userID := c.Query("user_id")

	var review models.Review
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).Find(&review).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"error getting review from database"})
		return
	}

	if err:=config.DB.Delete(&review).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"error":"Database error"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"Review deleted"})
}