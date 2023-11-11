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
	
	c.JSON(200, gin.H{"plans": plans})
	
	
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
		c.JSON(http.StatusBadRequest,gin.H{"error":"error binding"})
		return
	}
	validate.RegisterValidation("plan",ValidatePlan)
	if err := validate.Struct(user); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Please select correct plan"})
		return
	}

	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	userIDString:=fmt.Sprint(userID)
	fmt.Println(userIDString)
		if err := config.Client.Set(ctx,"plan"+userIDString,user.Plan,30*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
			return
		}
	c.JSON(200,gin.H{"message":"ok"})
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
			"error":"User not found",
		})
		return
	}
	var existingMember models.Membership

	if err:=config.DB.Where("user_id = ?",user.UserID).First(&existingMember).Error; err == nil {
			c.JSON(http.StatusConflict,gin.H{"error":"Already have membership"})
			return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadGateway,gin.H{"error":"database error"})
		return
	}
	
	userPlan,err:=config.Client.Get(ctx,"plan"+userID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error":"plan not found"})
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
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid plan"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}


	userPlan,err:=config.Client.Get(ctx,"plan"+userIDString).Result()
	if err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error":"otp not found"})
		return
	}
	var d time.Duration
	if userPlan=="5M"{
		d=time.Hour*24*150
	}else if userPlan=="1Y"{
		d=time.Hour*24*365
	}else if userPlan=="3Y"{
		d=time.Hour*24*1095
	}else {
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid plan"})
	}


	var user models.Membership

	user.UserID=uint64(UserID)
	user.RazorpaySubscriptionID=paymentID
	user.Plan=userPlan
	user.StartedOn=time.Now()
	user.ExpiresAt=time.Now().Add(d)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"database error"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"status": true})
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


//InvoiceDownload helps to download the invoice
func InvoiceDownload(c *gin.Context) {
	content, err := generateInvoice()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate invoice"})
			return
		}

		// Set headers for downloading the PDF
		c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
		c.Header("Content-Type", "application/pdf")
		c.Data(http.StatusOK, "application/pdf", content)
}


func generateInvoice() ([]byte,error) {
	var buf []byte

	return buf,nil
}



