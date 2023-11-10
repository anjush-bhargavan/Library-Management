package controllers

import (
	"net/http"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)


//ViewUser handles admin to view a user details
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


//ViewUsers function shows all users
func ViewUsers(c *gin.Context) {
	var users []models.User

	config.DB.Find(&users)

	c.JSON(http.StatusOK,users)
}


//UpdateUser handles the admin to update user details
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
	if err := validate.Struct(user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Please fill all fields"})
		return
	}

	config.DB.Save(&user)
	c.JSON(http.StatusOK,user)
}

//DeleteUser handles admin to delete user by id
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