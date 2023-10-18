package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

func ViewUser(c *gin.Context) {
	id :=c.Param("id")
	var user models.User

	if err :=config.DB.First(&user,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"User not found",
		})
		return
	}

	c.JSON(http.StatusOK,user)
}

func AddUser(c *gin.Context) {
	var user models.User

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Binding error"})
		return
	}

	config.DB.Create(&user)
	c.JSON(http.StatusOK,gin.H{"message":"User created succesfully"})
}

func UpdateUser(c *gin.Context) {
	id :=c.Param("id")
	var user models.User

	if err :=config.DB.First(&user,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"User not found",
		})
		return
	}

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Binding error"})
		return
	}

	config.DB.Save(&user)
	c.JSON(http.StatusOK,user)
}

func DeleteUser(c *gin.Context) {
	id :=c.Param("id")
	var user models.User

	if err :=config.DB.First(&user,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"User not found",
		})
		return
	}

	config.DB.Delete(&user)
	c.JSON(http.StatusOK,gin.H{"message":"User deleted"})

}