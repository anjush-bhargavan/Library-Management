package controllers

import (
	"net/http"
	"github.com/anjush-bhargavan/library-management/auth"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


func UserLoginPage(c *gin.Context){
	c.JSON(200,gin.H{
		"message":"Login with email & password",
	})
}

func UserLogin(c *gin.Context){
	type login struct {
		Email     string    `json:"email"`
		Password  string	`json:"password"`
	}

	var temp login

	if err :=c.ShouldBindJSON(&temp);err != nil {
		 c.JSON(http.StatusBadGateway,gin.H{
			"error":"Binding error",
		 })
	}

	var user models.User

	if err := config.DB.Where("email = ?",temp.Email).First(&user).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"User not found"})
		return
	}


	if err :=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(temp.Password)); err == nil{
		token,err:=auth.GenerateToken(user.Email,user.Role)
		if err != nil{
			c.JSON(http.StatusBadGateway,gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(200,gin.H{
			"token":token,
		})

		
		c.Header("Authorization","Bearer "+token)
		c.Redirect(http.StatusSeeOther,"/user/home")
	}else{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid password"})
	}


}

func  UserSignupPage(c *gin.Context){
	c.JSON(200,gin.H{
		"message":"Sign up to Continue",
	})
}

func UserSignup(c *gin.Context){
	var user models.User

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"error" : "Binding error",
		})
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

	c.JSON(200,gin.H{"message":"Signup successful"})

}

func HomePage(c *gin.Context){
	data,_:=c.Get("email")
	email:=data.(string)
	c.JSON(200,gin.H{
		"message":"Welcome to homepage "+email,
	})
}

func UserLogout(c *gin.Context) {
	c.Redirect(http.StatusSeeOther,"/login")
}