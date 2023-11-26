package controllers

import (
	"net/http"
	"time"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//InholdBook function shows the details of user holding book
func InholdBook(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var bookout models.BooksOut
	if err := config.DB.Where("user_id = ?",userID).First(&bookout).Error; err == gorm.ErrRecordNotFound {
		
		c.JSON(http.StatusAccepted,gin.H{"status":"Success",
										"message":"No books in hand",
										"data":nil,
										})
		return
		
	}else if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Database error",
										"data":err.Error(),
										})
		return
	}
	var book models.Book
	if err := config.DB.Where("id = ?",bookout.BookID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
											})
		return
	}

	c.JSON(200,gin.H{"status":"Success",
					"message":"Inhold Book details ",
					"data":book,
					})
}


//ReturnBook fucntion place a order to collect the book
func ReturnBook(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var bookout models.BooksOut
	if err := config.DB.Where("user_id = ?",userID).First(&bookout).Error; err == gorm.ErrRecordNotFound {
		
		c.JSON(http.StatusAccepted,gin.H{"status":"Failed",
										"message":"No books in hand",
										"data":nil,
										})
		return
	}else if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Database error",
										"data":err.Error(),
										})
		return
	}

	var order models.Orders

	order.UserID=bookout.UserID
	order.BookID=bookout.BookID
	order.Type="return"
	order.OrderedOn=time.Now()

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
											})
		return
	}
	c.JSON(200,gin.H{"status":"Success",
					"message":"Return order placed",
					"data":order,
					})
	
}