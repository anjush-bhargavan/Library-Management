package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

//ViewFine handles users to view their fine
func ViewFine(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	var user models.FineList

	if err:=config.DB.Where("user_id = ?",userID).First(&user).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200,gin.H{	"status":"Success",
							"message":"You have no fine till now",
							"data":err.Error(),
						})
		return
	}else if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
											})
		return
	}

	c.JSON(200,gin.H{   "status":"Success",
						"message":"You have fine of : ",
						"data":fmt.Sprint(user.Fine),
					})
}


//PayFine handle helps users to pay the fine
func PayFine(c *gin.Context) {
	userID:=c.Query("id")
	var user models.User
	if err :=config.DB.First(&user,userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
											"status":"Failed",
											"message":"User not found",
											"data":err.Error(),
										})
		return
	}

	var fine models.FineList
	if err:=config.DB.Where("user_id = ?",userID).First(&fine).Error; err == gorm.ErrRecordNotFound {
		c.JSON(200,gin.H{"status":"Success",
						"message":"You have no fine till now",
						"data":"",
					})
		return
	}else if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}

	finelink:="fine"
	if err := config.Client.Set(ctx,"plan"+userID,finelink,30*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"status":"Failed",
													"message":"Error setting data in redis server",
													"data":err.Error(),
												})
		return
	}

	amountInPaisa:= fine.Fine *100

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
