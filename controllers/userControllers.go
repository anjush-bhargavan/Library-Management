package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/anjush-bhargavan/library-management/auth"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//UserLoginPage handles get login page
func UserLoginPage(c *gin.Context){
	c.JSON(200,gin.H{
		"message":"Login with email & password",
	})
}


//UserLogin handles login and create jwt token
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


//UserSignupPage handles get signup page
func  UserSignupPage(c *gin.Context){
	c.JSON(200,gin.H{
		"message":"Sign up to Continue",
	})
}


var ctx = context.Background()


//UserSignup handles post signup form and validation
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
	
	otp:=auth.GenerateOTP(6)
	auth.SendOTPByEmail(otp,user.Email)

	jsonData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":"Failed to marshal json data"})
		return
	}
	
	if err := config.Client.Set(ctx,"otp",otp,30*time.Second).Err(); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}
	if err := config.Client.Set(ctx,"user",jsonData,30*time.Second).Err(); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}
	
	
	fmt.Println(otp)


	c.JSON(http.StatusOK,gin.H{"message":"go to verfication page"})

}


//VerifyOTPPage shows the verify otp
func VerifyOTPPage(c *gin.Context) {
	c.JSON(200,gin.H{"message":"verify otp"})
}


//VerifyOTP handles verifying otp and saving user data in database
func VerifyOTP(c *gin.Context) {
	var userData models.User
	type OTPString struct{
		Otp string `json:"otp"`
	}
	var user OTPString 
	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"error" : "Binding error",
		})
		return
	}
	if user.Otp ==""{
		c.JSON(http.StatusNotFound,gin.H{"error":"otp not entered"})
		return
	}
	otp,err:=config.Client.Get(ctx,"otp").Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error":"otp not found"})
		return
	}
	if auth.ValidateOTP(otp,user.Otp) {
		user,err := config.Client.Get(ctx,"user").Result()
		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{"error":"user details missing"})
		}
		err = json.Unmarshal([]byte(user),&userData)
		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{"error":"error in unmarshaling json data"})
		}
		config.DB.Create(&userData)
		c.JSON(http.StatusOK,gin.H{"message":"signup successful"})
	}
}


//HomePage handles get homepage
func HomePage(c *gin.Context){
	data,_:=c.Get("email")
	email:=data.(string)
	c.JSON(200,gin.H{
		"message":"Welcome to homepage "+email,
	})
}


//UserLogout handles logout
func UserLogout(c *gin.Context) {
	c.Redirect(http.StatusSeeOther,"/login")
}