package controllers

import (
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//ViewHistory function handles admin to view all history
func ViewHistory(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var history []models.History

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&history)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":history,
					})
}

//ViewBooksOut function shows the table of booksout
func ViewBooksOut(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var booksOut []models.BooksOut

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&booksOut)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":booksOut,
					})
}

//ViewOrders function handles admin to view all history
func ViewOrders(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var orders []models.Orders

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&orders)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":orders,
					})
}
