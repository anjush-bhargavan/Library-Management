package controllers

import (
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//ViewFeedbacks handles admin to view received feedback
func ViewFeedbacks(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var feedbacks []models.FeedBack

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&feedbacks)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Feedbscks fetched succesfully",
						"data":feedbacks,
					})
} 