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

	var history models.History
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).First(&history).Error; err != nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"You didn't borrowed this book",
											"data":err.Error(),
										})
		return
	}else if err == nil && history.Status =="pending"{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"You didn't borrowed this book",
											"data":err.Error(),
										})
		return
	}

	var review models.Review

	var existingReview models.Review
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).Find(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Review already exists",
											"data":existingReview,
										})
		return
	}else if err !=gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Error getting review from databse",
											"data":err.Error(),
										})
		return
	}

	if err :=c.ShouldBindJSON(&review);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Error Binding",
											"data":err.Error(),
										})
		return
	}
	review.UserID=userID
	review.BookID=uint64(bookID)
	
	if err := config.DB.Create(&review).Error; err!= nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	
	var book models.Book
	if err := config.DB.Where("book_id = ?",bookID).First(&book).Error;err != nil{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	book.Rating=(book.Rating+review.Rating)/2
	if err := config.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Review  added succesfully",
						"data":review,
					})
}


//ShowReview function shows users review of books
func ShowReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)

	var reviews []models.Review

	if err :=config.DB.Find(&reviews).Where("id = ?",bookID).Error; err != nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Error getting reviews of book",
											"data":err.Error(),
										})
		return
	}
	var sum uint64
	for _,review:=range reviews{
		sum+=review.Rating
	}
	Rating:=int(sum)/len(reviews)

	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message1":"Rating",
								"data1":Rating,
								"message2":"count of reviews",
								"data2":len(reviews),
								"massage3":"Reviews of book",
								"data3":reviews,

							})
}


//EditReview function helps users to edit review
func EditReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var review models.Review
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).Find(&review).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Error getting review from databse",
											"data":err.Error(),
										})
		return
	}

	if err :=c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}


	if err:=config.DB.Save(&review).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Review ",
								"data":review,
							})
}


//DeleteReview handles users to delete their review
func DeleteReview(c *gin.Context) {
	id:=c.Param("id")
	bookID,_:=strconv.Atoi(id)
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var review models.Review
	if err := config.DB.Where("book_id = ? AND user_id = ?",bookID,userID).Find(&review).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Error getting review from database",
											"data":err.Error(),
										})
		return
	}

	if err:=config.DB.Delete(&review).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"Success",
								"message":"Review Deleted",
								"data":nil,
							})
}