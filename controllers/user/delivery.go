package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ctx =context.Background()

//DeliveryDetails show the details regarding the delivery
func DeliveryDetails(c *gin.Context) {
	id :=c.Param("id")
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	userIDString:=fmt.Sprint(userID)

	if err := config.DB.Model(models.Membership{}).Where("user_id = ? AND is_active = ?",userID,true).Error; err != nil {
		c.JSON(http.StatusUnauthorized,gin.H{"status":"Failed",
											"message":"You have to take membership before borrowing books",
											"data":err.Error(),
										})
	return
	}
	bookID,_ :=strconv.Atoi(id)
	bookIDUint := uint64(bookID)

	var bookout models.BooksOut

	if err :=config.DB.Where("book_id = ?",bookIDUint).First(&bookout).Error; err == nil {
		c.JSON(http.StatusFound, gin.H{"status":"Failed",
										"message":"Book is not currently available, Will be available after ",
										"data":bookout.ReturnDate,
										})
		return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Book not found",
										"data":err.Error(),
										})
		return
	}

	var history models.History

	if err :=config.DB.Where("book_id = ? AND status = ?",bookIDUint,"pending").First(&history).Error; err == nil {
		c.JSON(http.StatusFound, gin.H{"status":"Failed",
										"message":"Book is currently ordered by another member",
										"data":nil,
										})
		return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Book not found",
										"data":err.Error(),
										})
		return
	}


	Date := time.Now().AddDate(0,0,5).Format("01-02-2008")
	c.JSON(200,gin.H{"status":"Success",
					"message":"Available date for delivery",
					"data":Date,
					})
	var existingBookid uint64

	if err:=config.DB.Model(&models.BooksOut{}).Where("user_id = ?",userID).Pluck("BookID",&existingBookid).Error;err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	fmt.Println(existingBookid)
	if existingBookid==0{
		c.JSON(200,gin.H{"status":"Success",
						"message":"No books in hand",
						"data":"",
					})
	}else{
		var existingBook models.Book

		if err :=config.DB.First(&existingBook,existingBookid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
											"message":"Book not found",
											"data":err.Error(),
											})
			return
		}
		c.JSON(200,gin.H{"status":"Success",
						"message":"You have to return this book during delivery",
						"data":existingBook,
					})
	}
		

		var book models.Book
		
		if err :=config.DB.First(&book,bookIDUint).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
											"message":"Book not found",
											"data":err.Error(),
											})
			return
		}
		
		if err := config.Client.Set(ctx,"bookid"+userIDString,bookIDUint,30*time.Minute).Err(); err != nil {
			c.JSON(400,gin.H{"status":"Failed",
							"message":"Failed to set bookID in redis",
							"data": err.Error(),
							})
			return
		}

	c.JSON(http.StatusOK,book)
}

//Delivery adds the order to history
func Delivery(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	userIDString:=fmt.Sprint(userID)

	var history models.History

	bookIDString,err:=config.Client.Get(ctx,"bookid"+userIDString).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"status":"Failed",
											"message":"BookID not found in redis",
											"data":err,
										})
		return
	}
	bookID, _ := strconv.ParseUint(bookIDString, 0, 64)

	history.BookID=bookID
	history.UserID=userID
	history.OrderedOn=time.Now()

	var user models.History

	if err:=config.DB.Where("user_id = ? AND status = ?",userID,"pending").First(&user).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"status":"Failed",
											"message":"Already one order is pending",
											"data":"",
										})
		return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err,
										})
		return
	}

	
	if err :=config.DB.Create(&history).Error;err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err,
										})
		return
	}

	c.JSON(200,gin.H{"status":"Success",
					"message":"Order placed, keep the other book ready during delivery",
					"data":history,
					})
}

//CancelOrder function handles users to cancel their order
func CancelOrder(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var user models.History

	if err:=config.DB.Where("user_id = ? AND status = ?",userID,"pending").First(&user).Error; err != nil {
		c.JSON(http.StatusConflict,gin.H{"status":"Failed",
										"message":"No order is pending",
										"data":err,
										})
		return
	}
	user.Status="cancelled"
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error while saving",
											"data":err,
										})
		return
	}
	c.JSON(200,gin.H{"status":"Success",
					"message":"Order cancelled",
					"data":"",
				})
}