package controllers

import (
	"net/http"
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//AddReview function helps users to add review
func AddReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var review models.Review

	var existingReview models.Review
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).Find(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Review already exists"})
		return
	}else if err !=gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadGateway,gin.H{"error":"error getting review from database"})
		return
	}

	if err :=c.ShouldBindJSON(&review);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"error binding"})
		return
	}
	review.UserID=userID
	review.BookID=uint64(bookID)
	
	if err := config.DB.Create(&review).Error; err!= nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Database error"})
		return
	}
	c.JSON(200,review)
}


//ShowReview function shows users review of books
func ShowReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)

	var reviews []models.Review

	if err :=config.DB.Find(&reviews).Where("id = ?",bookID).Error; err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":"error in getting reviews"})
		return
	}
	var sum uint64
	for _,review:=range reviews{
		sum+=review.Rating
	}
	Rating:=int(sum)/len(reviews)

	c.JSON(http.StatusOK,gin.H{"Rating":Rating,"Review":reviews})
}


//EditReview function helps users to edit review
func EditReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var review models.Review
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).Find(&review).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"error getting review from database"})
		return
	}

	if err :=c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Binding error"})
		return
	}


	if err:=config.DB.Save(&review).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"error":"Database error"})
		return
	}
	c.JSON(http.StatusOK,review)
}


//DeleteReview handles users to delete their review
func DeleteReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

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