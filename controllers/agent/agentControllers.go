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
	var orders []models.Orders

	if err:= config.DB.Where("status = ?","pending").Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Database error getting orders",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Orders fetched succesfully",
						"data":orders,
					})
}

//GetOrder function shows agent the order details using order_id
func GetOrder(c *gin.Context) {
	orderID :=c.Param("id")

	var order models.Orders
	
	if err := config.DB.Where("id = ?",orderID).First(&order).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error getting user",
											"data":err.Error(),
										})
		return
	}

	var result struct{
		ID			uint64
		UserName	string
		Phone		string
		Address		string
		BookName	string
		Type 		string
		OrderedOn	time.Time
		Status		string
	}
	query := `
    SELECT 
        orders.id,
        users.user_name,
        users.phone,
        users.address,
        books.book_name,
		orders.type,
        orders.ordered_on,
        orders.status
    FROM 
        orders
    JOIN 
        users ON orders.user_id = users.user_id
    JOIN 
        books ON orders.book_id = books.id
    WHERE 
        orders.id = ?`

	if err := config.DB.Raw(query, order.ID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{	"status":"Failed",
												"message":"Database error while joining",
												"data":err.Error(),
											})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Order details",
						"data":result,
						})
}


//UpdateOrders fuction handles the update on the history table by agent using the user id 
func UpdateOrders(c *gin.Context) {
	type Delivery struct{
		UserID 		uint64	`json:"user_id"`
		Type        string	`json:"type"`
		Status 		string	`json:"status"`
	}
	var first Delivery
	if err := c.ShouldBindJSON(&first); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
	}

	
	 
	var Currentorder models.Orders

	if err := config.DB.Where("user_id = ? AND status = ? AND type = ?",first.UserID,"pending",first.Type).First(&Currentorder).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	Currentorder.Status="ok"

	if err := config.DB.Save(&Currentorder).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error while saving",
											"data":err.Error(),
										})
		return
	}

	if Currentorder.Type=="delivery"{
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
		var history models.History

		history.UserID=Currentorder.UserID
		history.BookID=Currentorder.BookID
		history.RentedOn=time.Now()

		if err := config.DB.Create(&history).Error ; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Database error while creating history",
												"data":err.Error(),
											})
			return
		}
	}
	if Currentorder.Type=="return"{
		var history models.History
		if err :=config.DB.Where("user_id = ? AND book_id = ?",Currentorder.UserID,Currentorder.BookID).First(&history).Error; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
												})	
			return					
		}
		history.Status="returned"
		history.ReturnedOn=time.Now()
		if err := config.DB.Save(&history).Error; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
												})	
		return		
		}

		var book models.Book
		if err := config.DB.Where("id = ?",Currentorder.BookID).First(&book).Error;err != nil{
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

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Order updated succesfully",
						"data":Currentorder,
					})
}


