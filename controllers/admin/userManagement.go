package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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


//AddUser handles the admin to add user to database
func AddUser(c *gin.Context) {
	var user models.User

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"Binding error"})
		return
	}
	var existingUser models.User
	if err := config.DB.Where("email = ?",user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Email already in use"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}

	if err := config.DB.Where("phone = ?",user.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Phone number already in use"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{ "error":"Failed to hash password",
	})
	return
	}
	user.Password=string(hashedPassword)

	config.DB.Create(&user)
	c.JSON(http.StatusOK,gin.H{"message":"User created succesfully"})
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