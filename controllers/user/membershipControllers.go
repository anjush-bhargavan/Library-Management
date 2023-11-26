package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jung-kurt/gofpdf"
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
											"data":nil,
										})
			return
	}else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest,gin.H{
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
	userID :=c.Query("user_id")
	fmt.Println(pID)
	fmt.Println("Fully successful")

	c.HTML(http.StatusOK, "success.html", gin.H{
		"paymentID": pID,
		"userID":userID,
	})
}

//InvoiceDownload handles the user to download the invoice of membership
func InvoiceDownload(c *gin.Context) {
	userIDQuery :=c.Query("user_id")
	userIDint,_:=strconv.Atoi(userIDQuery)
	userID:=uint64(userIDint)
	content, err := generateMembershipSubscription(userID)
	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
											"status":"Failed",
											"message":"Failed to generate membership subscription PDF",
											"data":err.Error(),
										})
		return
	}

	// Set headers for downloading the PDF
	c.Header("Content-Disposition", "attachment; filename=membership_subscription.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", content)
	
}

//generateMembershipSubscription function finds the details of member and generate invoice
func generateMembershipSubscription(userID uint64) ([]byte, error) {
    var result struct {
        UserID                uint64
        UserName              string
        Phone                 string
        Email                 string
        Plan                  string
        RazorpaySubscriptionID string
        StartedOn             time.Time
        ExpiresAt             time.Time
    }

    query := `
    SELECT 
        users.user_id,
        users.user_name,
        users.phone,
        users.email,
        memberships.plan,
        memberships.razorpay_subscription_id,
        memberships.started_on,
        memberships.expires_at
    FROM 
        memberships
    JOIN 
        users ON memberships.user_id = users.user_id
    WHERE 
        memberships.user_id = ?`

    if err := config.DB.Raw(query, userID).Scan(&result).Error; err != nil {
        log.Println("error getting records from database")
        return nil, err
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
  
    pdf.AddPage()

	pdf.SetFont("Times", "", 20)


	pdf.SetFillColor(135, 206, 250) 
	pdf.CellFormat(190, 10, "Golib.Online", "1", 0, "C", true, 0, "")


	pdf.SetFont("Helvetica", "", 12)

	
	pdf.Ln(15)
	pdf.SetFillColor(180, 200, 200) 
	pdf.CellFormat(190, 10, "Thank you for taking the membership!", "1", 0, "C", true, 0, "")
	pdf.Ln(15)

	
	pdf.SetFillColor(255, 255, 255) 


	pdf.Cell(0, 10, fmt.Sprintf("User ID: %d", result.UserID))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("User Name: %s", result.UserName))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Phone: %s", result.Phone))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Email: %s", result.Email))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Plan: %s", result.Plan))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Razorpay Subscription ID: %s", result.RazorpaySubscriptionID))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Started On: %s", result.StartedOn.Format("01-Jan-2006")))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Expiration Date: %s", result.ExpiresAt.Format("01-Jan-2006")))
    pdf.Ln(10)

    pdf.Cell(0, 10, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
    pdf.Ln(10)

    pdf.Ln(10)
	pdf.MultiCell(0, 8, "You can now borrow books from Golib.Online and keep looking forward to an enriching reading experience.", "0", "L", false)
    pdf.Ln(10)

    if err := pdf.OutputFileAndClose("membership_subscription.pdf"); err != nil {
        log.Println("error making pdf")
        return nil, err
    }

    content, err := os.ReadFile("membership_subscription.pdf")
    if err != nil {
        log.Println("error reading pdf file")
        return nil, err
    }

    return content, nil
}



