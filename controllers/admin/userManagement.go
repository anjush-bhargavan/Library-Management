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
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
											})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"User fetched succesfully",
						"data":user,
					})
}


//ViewUsers function shows all users
func ViewUsers(c *gin.Context) {
	var users []models.User

	config.DB.Find(&users)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Users fetched succesfully",
						"data":users,
					})
}


//UpdateUser handles the admin to update user details
func UpdateUser(c *gin.Context) {
	id :=c.Param("id")
	var user models.User

	if err :=config.DB.First(&user,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
										})
		return
	}

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}
	if err := validate.Struct(user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	config.DB.Save(&user)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"User updated succesfully",
						"data":user,
					})
}

//DeleteUser handles admin to delete user by id
func DeleteUser(c *gin.Context) {
	id :=c.Param("id")
	var user models.User

	if err :=config.DB.First(&user,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
										})
		return
	}

	config.DB.Delete(&user)
	c.JSON(http.StatusOK,gin.H{	"status":"Failed",
								"message":"User deleted",
								"data":nil,
							})

}