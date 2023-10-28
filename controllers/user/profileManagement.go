package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//UserProfile handles to get profile page of user
func UserProfile(c *gin.Context) {
	data,_:=c.Get("email")
	email:=data.(string)
	var user models.User

	if err :=config.DB.Where("email = ?",email).Omit("user_id","password","role").First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"User not found",
		})
		return
	}
	c.JSON(http.StatusOK,user)
}

//ProfileUpdate handles the updates of userprofile
func ProfileUpdate(c *gin.Context) {
	data,_:=c.Get("email")
	email:=data.(string)
	var user models.User

	if err :=config.DB.Where("email = ?",email).Omit("user_id","password","role").First(&user).Error; err != nil {
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

//Membership handles the membership of users
func Membership(c *gin.Context) {

}

//ViewHistory handles to show the history of book taken by user
func ViewHistory(c *gin.Context) {
	
}