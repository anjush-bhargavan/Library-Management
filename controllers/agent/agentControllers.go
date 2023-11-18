package controllers

import (
	"net/http"
	"time"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//ViewOrders function handles the agent to view the orders that are pending
func ViewOrders(c *gin.Context) {
	var orders []models.History

	if err:= config.DB.Where("status = ?","pending").Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Database error getting orders",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,orders)
}

//GetOrder function shows agent the order details using order_id
func GetOrder(c *gin.Context) {
	orderID :=c.Param("id")

	var order models.History
	
	if err := config.DB.Where("id = ?",orderID).First(&order).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error getting user",
											"data":err.Error(),
										})
		return
	}

	var result struct{
		ID			uint64
		FirstName	string
		Phone		string
		Address		string
		BookName	string
		OrderedOn	time.Time
		Status		string
	}
	query := `
    SELECT 
        histories.id,
        users.user_name,
        users.phone,
        users.address,
        books.book_name,
        histories.ordered_on,
        histories.status
    FROM 
        histories
    JOIN 
        users ON histories.user_id = users.user_id
    JOIN 
        books ON histories.book_id = books.id
    WHERE 
        histories.id = ?`

if err := config.DB.Raw(query, order.ID).Scan(&result).Error; err != nil {
    c.JSON(http.StatusBadGateway, gin.H{	"status":"Failed",
											"message":"Database error while joining",
											"data":err.Error(),
										})
    return
}

	c.JSON(200,result)
}


//UpdateOrders fuction handles the update on the history table by agent using the user id 
func UpdateOrders(c *gin.Context) {
	type Order struct{
		UserID 		uint64	`json:"user_id"`
		Status 		string	`json:"status"`
	}
	var first Order
	if err := c.ShouldBindJSON(&first); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
	}

	
	
	var Currentorder models.History

	if err := config.DB.Where("user_id = ? AND status = ?",first.UserID,"pending").First(&Currentorder).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	Currentorder.Status=first.Status

	if err := config.DB.Save(&Currentorder).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error while saving",
											"data":err.Error(),
										})
		return
	}
	if first.Status=="delivered"{
		var bookout models.BooksOut
		bookout.UserID=Currentorder.UserID
		bookout.BookID=Currentorder.BookID
		bookout.OutDate=time.Now()
		bookout.ReturnDate=time.Now().AddDate(0,0,20)
		if err :=config.DB.Create(&bookout).Error;err != nil{
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error in bookout",
												"data":err.Error(),
											})
			return
		}
	}
	if first.Status=="returned"{

		var book models.Book
		if err := config.DB.Where("book_id = ?",Currentorder.BookID).First(&book).Error;err != nil{
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
											})
			return
		}
		book.OrderCount++
		if err := config.DB.Save(&book).Error; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
											})
			return
		}
		
		var bookOut models.BooksOut
		if err := config.DB.Where("user_id = ?",Currentorder.UserID).Find(&bookOut).Error; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error in getting user",
												"data":err.Error(),
											})
			return
		}
		today:=time.Now()
		if bookOut.ReturnDate.Before(today) {
			var fineduser models.FineList 
			fineduser.UserID=Currentorder.UserID
			duration:=time.Since(bookOut.ReturnDate)
			fineduser.Fine=  uint64((duration.Hours() / 24 )* 10)

			if err :=config.DB.Create(&fineduser).Error; err!=nil {
				c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
													"message":"Database error",
													"data":err.Error(),
												})
				return
			}
		}
	
		if err:=config.DB.Delete(&bookOut).Error; err != nil{
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
											})
			return
		}
		c.JSON(http.StatusOK,gin.H{	"status":"Success",
									"message":"Books out table updated",
									"data":nil,
								})
	}

	c.JSON(200,Currentorder)
}


