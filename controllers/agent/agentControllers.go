package controllers

import (
	"net/http"
	"time"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//ViewOrders function handles the agent to view the orders
func ViewOrders(c *gin.Context) {
	var orders []models.History

	if err:= config.DB.Where("status = ?","pending").Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":"Database error"})
		return
	}

	c.JSON(200,orders)
}

//GetOrder function shows agent the order details
func GetOrder(c *gin.Context) {
	orderID :=c.Param("id")

	var order models.History
	
	if err := config.DB.Where("id = ?",orderID).First(&order).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"error": "database error"})
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
    c.JSON(http.StatusBadGateway, gin.H{"error": "database error while joining"})
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
		c.JSON(http.StatusBadGateway,gin.H{"error":"Binding error"})
	}
	
	var Currentorder models.History

	if err := config.DB.Where("user_id = ? AND status = ?",first.UserID,"pending").First(&Currentorder).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"erro":"Database error"})
		return
	}
	Currentorder.Status=first.Status

	if err := config.DB.Save(&Currentorder).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Database error while saving"})
		return
	}

	c.JSON(200,Currentorder)
}