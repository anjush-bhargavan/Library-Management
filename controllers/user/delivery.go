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

	if err := config.DB.Model(models.Membership{}).Where("user_id = ?",userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized,gin.H{"error":"You have to take membership before borrowing books"})
		return
	}


	Date := time.Now().AddDate(0,0,5)
	c.JSON(200,gin.H{"Available date for delivery":Date})
	var existingBookid uint64

	if err:=config.DB.Model(&models.BooksOut{}).Where("user_id = ?",userID).Pluck("BookID",&existingBookid).Error;err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"database error"})
		return
	}
	fmt.Println(existingBookid)
	if existingBookid==0{
		c.JSON(200,gin.H{"message":"no books in hand"})
	}else{
		var existingBook models.Book

		if err :=config.DB.First(&existingBook,existingBookid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":"book not found",
			})
			return
		}
		c.JSON(200,gin.H{"You have to return this book during delivery":existingBook})
	}
		

		var book models.Book
		bookID,_ :=strconv.Atoi(id)
		bookIDString := uint64(bookID)
		if err :=config.DB.First(&book,bookIDString).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":"book not found",
			})
			return
		}
		
		if err := config.Client.Set(ctx,"bookid"+userIDString,bookIDString,30*time.Minute).Err(); err != nil {
			c.JSON(400,gin.H{"error": err.Error()})
			return
		}

		// if err := config.Client.Set(ctx,"bookid"+userIDString,bookIDString,30*time.Minute).Err(); err != nil {
		// 	c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		// 	return
		// }

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
		c.JSON(http.StatusNotFound,gin.H{"error":"bookid not found"})
		return
	}
	bookID, _ := strconv.ParseUint(bookIDString, 0, 64)

	history.BookID=bookID
	history.UserID=userID
	history.OrderedOn=time.Now()

	var user models.History

	if err:=config.DB.Where("user_id = ? AND status = ?",userID,"pending").First(&user).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"erro":"Already one order is pending"})
		return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadGateway,gin.H{"error":"database error"})
		return
	}

	
	if err :=config.DB.Create(&history).Error;err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"error":"database error"+err.Error()})
		return
	}
	c.JSON(200,gin.H{"message":"order placed, keep the other book ready during delivery"})
}

//CancelOrder function handles users to cancel their order
func CancelOrder(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var user models.History

	if err:=config.DB.Where("user_id = ? AND status = ?",userID,"pending").First(&user).Error; err != nil {
		c.JSON(http.StatusConflict,gin.H{"erro":"No order is pending"})
		return
	}
	user.Status="cancelled"
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Database error while saving"})
		return
	}
}