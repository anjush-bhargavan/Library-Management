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
	"github.com/go-playground/validator/v10"
)

var validate =validator.New()



//UserLogin handles login and create jwt token
func UserLogin(c *gin.Context){
	type login struct {
		Email     string    `json:"email" validate:"required"`
		Password  string	`json:"password" validate:"required"`
	}

	var temp login

	if err :=c.ShouldBindJSON(&temp);err != nil {
		 c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Binding error",
												"data":err.Error(),
											})
	}

	if err := validate.Struct(temp); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?",temp.Email).First(&user).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
										})
		return
	}

	if err :=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(temp.Password)); err == nil{
		token,err:=auth.GenerateToken(user.UserID,user.Email,user.Role)
		if err != nil{
			c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
												"message":"Failed to generate token",
												"data":err.Error(),
											})
			return
		}
		c.JSON(200,gin.H{	"status":"Success",
							"message":"Jwt token",
							"data":token,
							})

		
		c.Header("Authorization","Bearer "+token)
	}else{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Invalid password",
											"data":err.Error(),
										})
	}


}



var ctx = context.Background()


//UserSignup handles post signup form and validation
func UserSignup(c *gin.Context){
	var user models.User

	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
											})
		return
	}

	if err := validate.Struct(user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all the mandatory fields",
											"data":err.Error(),
										})
		return
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?",user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Email already in use",
											"data":"Choose another email",
										})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{	"status":"Failed",
														"message":"Database error",
														"data":err.Error(),
													})
		return
	}

	if err := config.DB.Where("phone = ?",user.Phone).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Phone number already in use",
											"data":err.Error(),
										})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{	"status":"Failed",
														"message":"Database error",
														"data":err.Error(),
													})
		return
	}

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{	"status":"Failed",
														"message":"Failed to hash password",
														"data":err.Error(),
													})
	return
	}
	user.Password=string(hashedPassword)
	user.UserName=user.FirstName+" "+user.LastName
	
	otp:=auth.GenerateOTP(6)
	auth.SendOTPByEmail(otp,user.Email)

	jsonData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Failed to marshal json data",
											"data":err.Error(),
										})
		return
	}
	
	if err := config.Client.Set(ctx,"otp"+user.Email,otp,30*time.Second).Err(); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{	"status":"Failed",
														"message":"Redis error",
														"data":err.Error(),
													})
		return
	}
	
	if err := config.Client.Set(ctx,"user"+user.Email,jsonData,30*time.Second).Err(); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"status":"Failed",
													"message":"Redis error",
													"data":err.Error(),
												})
		return
	}
	
	
	fmt.Println(otp)


	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Go to verfication page",
								"data":nil,
							})

}




//VerifyOTP handles verifying otp and saving user data in database
func VerifyOTP(c *gin.Context) {
	var userData models.User
	type OTPString struct{
		Email	string	`json:"email"`
		Otp 	string	`json:"otp"`
	}
	var user OTPString 
	if err :=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}
	if user.Otp ==""{
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"OTP not entered",
											"data":nil,
										})
		return
	}

	otp,err:=config.Client.Get(ctx,"otp"+user.Email).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"otp not found",
											"data":err.Error(),
										})
		return
	}
	if auth.ValidateOTP(otp,user.Otp) {
		user,err := config.Client.Get(ctx,"user"+user.Email).Result()
		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
												"message":"User details missing",
												"data":err.Error(),
											})
			return
		}
		err = json.Unmarshal([]byte(user),&userData)
		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
												"message":"Error in unmarshaling json data",
												"data":err.Error(),
											})
			return
		}
		config.DB.Create(&userData)
		c.JSON(http.StatusOK,gin.H{	"status":"Success",
									"message":"Signup successful",
									"data":userData,
								})
	}
}


//HomePage handles get homepage
func HomePage(c *gin.Context){
	data,_:=c.Get("email")
	email:=data.(string)

	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Welcome to homepage",
								"data":email,
							})
}


//IndexPage handles get indexpage
func IndexPage(c *gin.Context){
	

	c.HTML(http.StatusOK,"index.html",gin.H{	"status":"Success",
								"message":"Welcome to Golib",
								"data":nil,
							})
}

