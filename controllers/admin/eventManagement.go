package controllers

import (
	"net/http"
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//AddEvent handles admin to add any event
func AddEvent(c *gin.Context) {

	var event models.Event

	if err :=c.ShouldBindJSON(&event);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Error binding",
											"data":err.Error(),
										})
		return
	}

	if err :=config.DB.Create(&event).Error;err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Event added succesfully",
						"data":event,
					})
}


//EditEvent function helps admin to edit event
func EditEvent(c *gin.Context) {
	id:=c.Param("id")
	eventID,_:=strconv.Atoi(id)
	

	var event models.Event
	if err := config.DB.Where("event_id = ?",eventID).Find(&event).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Error getting event from database",
											"data":err.Error(),
										})
		return
	}

	if err :=c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}


	if err:=config.DB.Save(&event).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Event updated succesfully",
						"data":event,
					})
}

//ViewEvents handles admin to view events
func ViewEvents(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var events []models.Event

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&events)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"events fetched succesfully",
						"data":events,
					})
} 