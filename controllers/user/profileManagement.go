package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)



//UserProfile handles to get profile page of user
func UserProfile(c *gin.Context) {
	data,_:=c.Get("email")
	email:=data.(string)
	var user models.User

	if err :=config.DB.Where("email = ?",email).Omit("id","password","role").First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{  "status":"Failed",
											"message":"User not found",
											"data":err.Error(),
		
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Profile fetched succesfully",
						"data":user,
					})
}

//ProfileUpdate handles the updates of userprofile
func ProfileUpdate(c *gin.Context) {
	data,_:=c.Get("email")
	email:=data.(string)
	var user models.User

	if err :=config.DB.Where("email = ?",email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
										})
		return
	}
	user.UserName=user.FirstName+" "+user.LastName

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}


	config.DB.Save(&user)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Profile updated succesfully",
						"data":user,
					})
}


//ChangePassword function helps to change password
func ChangePassword(c *gin.Context) {
	data,_:=c.Get("email")
	email:=data.(string)
	var user models.User

	if err :=config.DB.Where("email = ?",email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
		
										})
		return
	}
	type password struct {
		Old  string `json:"old_password"`
		New  string `json:"new_password"`
		CNew string `json:"confirm_password"`
	}
	var newPAssword password
	if err :=c.ShouldBindJSON(&newPAssword); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}
	if err :=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(newPAssword.Old)); err != nil{
		c.JSON(http.StatusUnauthorized,gin.H{	"status":"Failed",
												"message":"Password not correct",
												"data":err.Error(),
											})
		return
	}
	if newPAssword.New==""{
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Password empty",
											"data":nil,
										})
		return
	}
	if newPAssword.New!=newPAssword.CNew{
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Password mismatch",
											"data":nil,
										})
		return
	}
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(newPAssword.New),bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{	"status":"Failed",
														"message":"Failed to hash password",
														"data":err.Error(),
													})
	return
	}
	user.Password=string(hashedPassword)
	config.DB.Save(&user)
	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Successfully changed password",
								"data":nil,
							})
	
}


//ViewHistory handles to show the history of book taken by user
func ViewHistory(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID := userIDContext.(uint64)

	var history []models.History
	if err :=config.DB.Find(&history).Where("id = ?",userID).Error; err != nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Error in getting history",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"History fetched succesfully",
						"data":history,
					})
}