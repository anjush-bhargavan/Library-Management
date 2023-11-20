package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//ViewEvents handles users to view upcoming events
func ViewEvents(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var events []models.Event

	offset := (page - 1)* pageSize

	currentDate := time.Now().Format("10-01-2008")
	if err :=config.DB.Where("date >= ?", currentDate).Order("id").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Upcoming events are :",
						"data":events,
					})
} 