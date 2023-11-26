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
	
	var member models.Membership
	if err := config.DB.Model(models.Membership{}).Where("user_id = ? AND is_active = ?",userID,true).First(&member).Error; err != nil {
		c.JSON(http.StatusUnauthorized,gin.H{"status":"Failed",
											"message":"You have to take membership before borrowing books",
											"data":err.Error(),
										})
	return
	}
	var fine models.FineList
	if err := config.DB.Where("user_id = ?",userID).First(&fine).Error; err == nil {
		c.JSON(http.StatusUnauthorized,gin.H{"status":"Failed",
											"message":"You have to pay fine before borrowing books, Your fine is :",
											"data":fine.Fine,
										})
		return
	}else if err != gorm.ErrRecordNotFound{
		c.JSON(http.StatusUnauthorized,gin.H{	"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
											})
		return
	}
	bookID,_ :=strconv.Atoi(id)
	bookIDUint := uint64(bookID)

	var book models.Book
		
	if err :=config.DB.First(&book,bookIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Book not found",
										"data":err.Error(),
										})
		return
	}

	var bookOutCount,pendingCount int64

	if err :=config.DB.Model(&models.BooksOut{}).Where("book_id = ?",bookIDUint).Count(&bookOutCount).Error; err!= nil{
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Error getting book out count",
										"data":err.Error(),
										})
	}
	if err :=config.DB.Model(&models.Orders{}).Where("book_id = ? AND status = ?",bookIDUint,"pending").Count(&pendingCount).Error; err!= nil{
		c.JSON(http.StatusNotFound, gin.H{"status":"Failed",
										"message":"Error getting  order count",
										"data":err.Error(),
										})
	}
	if bookOutCount+pendingCount >= int64(book.NoOfCopies) {
		c.JSON(http.StatusConflict, gin.H{"status":"Failed",
										"message":"Book is not currently available",
										"data":"Will be available after 30 days",
										})
		return
	}


	Date := time.Now().AddDate(0,0,5).Format("01-Jan-2006")
	c.JSON(200,gin.H{"status":"Success",
					"message":"Available date for delivery",
					"data":Date,
					})
	// var existingBookid uint64

	// Model(&models.BooksOut{}).Where("user_id = ?",userID).Pluck("BookID",&existingBookid)
	var existingBookOut models.BooksOut
	if err:=config.DB.Where("user_id = ?",userID).First(&existingBookOut).Error;err == nil{
		var existingBook models.Book

		if err :=config.DB.First(&existingBook,existingBookOut.BookID).Error; err != nil {
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
		
	}else if err == gorm.ErrRecordNotFound{
		c.JSON(200,gin.H{"status":"Success",
						"message":"No books in hand",
						"data":nil,
					})
	}else{
		c.JSON(http.StatusBadRequest,gin.H{"status":"Failed",
											"message":"Database error",
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

		c.JSON(200,gin.H{"status":"Success",
						"message":"Book",
						"data":book,
					})
}

//Delivery adds the order to history
func Delivery(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	userIDString:=fmt.Sprint(userID)

	var user models.Orders

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

	var order models.Orders

	bookIDString,err:=config.Client.Get(ctx,"bookid"+userIDString).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"status":"Failed",
											"message":"BookID not found in redis",
											"data":err,
										})
		return
	}
	bookID, _ := strconv.ParseUint(bookIDString, 0, 64)

	order.BookID=bookID
	order.UserID=userID
	order.Type="delivery"
	order.OrderedOn=time.Now()

	var bookout models.BooksOut
	if err := config.DB.Where("user_id = ?",userID).First(&bookout).Error; err == nil {
		var returnOrder models.Orders

		returnOrder.BookID=bookout.BookID
		returnOrder.UserID=bookout.UserID
		returnOrder.Type="return"
		returnOrder.OrderedOn=time.Now()

		if err :=config.DB.Create(&returnOrder).Error; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error",
												"data":err,
											})
			return
		}

		c.JSON(200,gin.H{	"status":"Success",
							"message":"Order placed, keep the other book ready during delivery",
							"data":returnOrder,
							})
	}else if err != gorm.ErrRecordNotFound{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err,
										})
		return
	}

	
	if err :=config.DB.Create(&order).Error;err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err,
										})
		return
	}

	c.JSON(200,gin.H{"status":"Success",
					"message":"Order placed, keep the other book ready during delivery",
					"data":order,
					})
}

//CancelOrder function handles users to cancel their order
func CancelOrder(c *gin.Context) {
	Type := c.Query("type")
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var user models.Orders

	if err:=config.DB.Where("user_id = ? AND status = ? AND type = ?",userID,"pending",Type).First(&user).Error; err != nil {
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