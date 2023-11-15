package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

//ShowPlans function shows the membership plans
func ShowPlans(c *gin.Context) {
	plans := []gin.H{
		{"plan": "5M", "description": "5 months you can use the platform to borrow books for 500"},
		{"plan": "1Y", "description": "1 year you can use the platform to borrow books for 1000"},
		{"plan": "3Y", "description": "3 years you can use the platform to borrow books for 2000"},
	}
	
	c.JSON(200, gin.H{
					"status":"Success",
					"message":"User Plans",
					"data":plans,
				})
}

var validate = validator.New()
// type Plan int 
// const (
// 	M5 Plan = iota
// 	Y1
// 	Y3

// )


//Subscription holds the detail of plan
type Subscription struct{
	Plan string `json:"plan" validate:"plan"`
}

//ValidatePlan function helps to validate plan details entered by user
func ValidatePlan(fl validator.FieldLevel) bool {
	plan := fl.Field().String()
	return plan == "3Y" || plan == "1Y" || plan == "5M"
}


//GetPlan function gets the plan from users
func GetPlan(c *gin.Context) {
	var user Subscription
	if err :=c.ShouldBindJSON(&user);err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
											"status":"Failed",
											"message":"Error binding",
											"data":err.Error(),
										})
		return
	}

	validate.RegisterValidation("plan",ValidatePlan)
	if err := validate.Struct(user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
											"status":"Failed",
											"message":"Please select correct plan",
											"data":err.Error(),
										})
		return
	}

	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	userIDString:=fmt.Sprint(userID)
	fmt.Println(userIDString)
		if err := config.Client.Set(ctx,"plan"+userIDString,user.Plan,30*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{
														"status":"Failed",
														"message":"Error seting data in redis",
														"data":err.Error(),
													})
			return
		}
	c.JSON(200,gin.H{
					"status":"Success",
					"message":"Plan selected",
					"data":user,
					})
}

type pageVariables struct {
	OrderID string
}

//Membership handles the membership of users
func Membership(c *gin.Context){
	userID:=c.Query("id")
	
	// userIDContext,_:=c.Get("user_id")
	
	// userID:=userIDContext.(uint64)
	// userIDString:=fmt.Sprint(userID)
	

	var user models.User

	if err :=config.DB.First(&user,userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
											"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
										})
		return
	}
	var existingMember models.Membership

	if err:=config.DB.Where("user_id = ? AND is_active = ?",user.UserID,true).First(&existingMember).Error; err == nil {
			c.JSON(http.StatusConflict,gin.H{
											"status":"Failed",
											"message":"Already have a membership",
											"data":err.Error(),
										})
			return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadGateway,gin.H{
											"status":"Failed",
											"message":"Databae error",
											"data":err.Error(),
										})
		return
	}
	
	userPlan,err:=config.Client.Get(ctx,"plan"+userID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{
											"status":"Failed",
											"message":"Plan not found",
											"data":err.Error(),
										})
		return
	}
	amount:=0
	if userPlan=="5M"{
		amount=500
	}else if userPlan=="1Y"{
		amount=1000
	}else if userPlan=="3Y"{
		amount=2000
	}else {
		c.JSON(http.StatusBadRequest,gin.H{
											"status":"Failed",
											"message":"Invalid plan",
											"data":err.Error(),
										})
	}
	amountInPaisa := amount * 100

	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))

	data := map[string]interface{}{
		"amount":   amountInPaisa,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)

	if err != nil {
		fmt.Printf("Problem getting repository information: %v\n", err)
		os.Exit(1)
	}
	
	value := body["id"]
	str := value.(string)


	homepageVariables := pageVariables{
		OrderID: str,
	}
	
	c.HTML(http.StatusOK, "app.html", gin.H{
		"userID":      userID,
		"totalPrice":  amountInPaisa / 100,
		"total":       amountInPaisa,
		"orderID":     homepageVariables.OrderID,
		"email":       user.Email,
		"phoneNumber": user.Phone,
	})
}

// RazorpaySuccess handles successful RazorPay payments.
func RazorpaySuccess(c *gin.Context) {
	userIDString := c.Query("user_id")
	UserID, _ := strconv.Atoi(userIDString)
	orderID := c.Query("order_id")
	paymentID := c.Query("payment_id")
	signature := c.Query("signature")
	paymentAmount := c.Query("total")
	amount, _ := strconv.Atoi(paymentAmount)
	fmt.Println(signature,orderID)
	rPay := models.RazorPay{
		UserID:          uint(UserID),
		RazorPaymentID:  paymentID,
		Signature:       signature,
		RazorPayOrderID: orderID,
		AmountPaid:      float64(amount),
	}

	result := config.DB.Create(&rPay)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
											"status":"Failed",
											"message":"Database error",
											"data":result.Error,
										})
		return
	}


	userPlan,err:=config.Client.Get(ctx,"plan"+userIDString).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{
											"status":"Failed",
											"message":"OTP Not found",
											"data":err.Error(),
										})
		return
	}

	if userPlan=="fine"{
		var user models.FineList


		if err := config.DB.Where("user_id = ?",UserID).Delete(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway,gin.H{
												"status":"Failed",
												"message":"Database error",
												"data":err.Error(),
											})
			return
		}
	}else{

		var d time.Duration
		if userPlan=="5M"{
			d=time.Hour*24*150
		}else if userPlan=="1Y"{
			d=time.Hour*24*365
		}else if userPlan=="3Y"{
			d=time.Hour*24*1095
		}else {
			c.JSON(http.StatusBadRequest,gin.H{
												"status":"Failed",
												"message":"Invalid plan",
												"data":err.Error(),
											})
		}

		var existingMember models.Membership

		if err:=config.DB.Where("user_id = ? AND is_active = ?",UserID,false).First(&existingMember).Error; err == nil {
			existingMember.RazorpaySubscriptionID=paymentID
			existingMember.Plan=userPlan
			existingMember.StartedOn=time.Now()
			existingMember.ExpiresAt=time.Now().Add(d)
			existingMember.IsActive=true

			if err := config.DB.Save(&existingMember).Error; err != nil {
				c.JSON(http.StatusBadGateway,gin.H{
													"status":"Failed",
													"message":"Database error",
													"data":err.Error(),
												})
				return
			}

		}else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadGateway,gin.H{
												"status":"Failed",
												"message":"Databae error",
												"data":err.Error(),
											})
			return
		}else{
			var user models.Membership

			user.UserID=uint64(UserID)
			user.RazorpaySubscriptionID=paymentID
			user.Plan=userPlan
			user.StartedOn=time.Now()
			user.ExpiresAt=time.Now().Add(d)

			if err := config.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusBadGateway,gin.H{
													"status":"Failed",
													"message":"Database error",
													"data":err.Error(),
												})
				return
			}

		}


		
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true})
}

// SuccessPage renders the success page.
func SuccessPage(c *gin.Context) {
	pID := c.Query("id")
	fmt.Println(pID)
	fmt.Println("Fully successful")

	c.HTML(http.StatusOK, "success.html", gin.H{
		"paymentID": pID,
	})
}




