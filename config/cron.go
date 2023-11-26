package config

import (
	"log"
	"net/smtp"
	"os"
	"time"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/robfig/cron/v3"
)

//InitCron function initializes the cron job and run the CheckMembership in each day
func InitCron() {
	c := cron.New()

	c.AddFunc("@daily", CheckMembership)
	c.AddFunc("@weekly",CheckBooksOut)

	c.Start()
}

//CheckMembership will check users membership is expired or not in each day
func CheckMembership() {
	var members []models.Membership
	now := time.Now()
	if err := DB.Where("expires_at < ? AND is_active = ?", now,true).Find(&members).Error; err != nil {
		log.Println("error finding records")
		return
	}


	for _, member := range members {
		var user models.User
		if err := DB.Where("user_id = ?",member.UserID).First(&user).Error; err != nil {
			log.Println("error finding records")
			return
		}
		expiresAt :=member.ExpiresAt.Format("01-Jan-2006")

		message := "Subject: Your GoLib Membership Has Expired\n\n" +
    "Dear " + user.UserName + ",\n\n" +
    "We hope this email finds you well. We wanted to inform you that your membership with GoLib has expired as of " + expiresAt + ".\n\n" +
    "Membership Details:\n" +
    "- Subscription ID: " + member.RazorpaySubscriptionID + "\n" +
    "- Member Name: " + user.UserName + "\n" +
    "- Plan: " + member.Plan + "\n" +
    "- Expiration Date: " + expiresAt + "\n\n" +
    "As your membership has expired, you may no longer have access to certain library services. To continue enjoying the benefits of GoLib, we encourage you to renew your membership at your earliest convenience.\n\n" +
    "Renew Membership\n\n" +
    "If you have already renewed your membership, please disregard this message.\n\n" +
    "Thank you for being a valued member of GoLib. If you have any questions or need assistance, feel free to contact our support team at [Support Email or Phone Number].\n\n" +
    "Best regards,\nThe GoLib Team"



		SendEmail(message,user.Email)
		member.IsActive = false
		DB.Save(&member)
	}

}

//CheckBooksOut function checks if returndate expired or not
func CheckBooksOut(){
	var booksout []models.BooksOut
	now := time.Now()
	if err := DB.Where("return_date < ?", now).Find(&booksout).Error; err != nil {
		log.Println("error finding records")
		return
	}

	for _, data := range booksout {

		var user models.User
		if err := DB.Where("user_id = ?",data.UserID).First(&user).Error; err != nil {
			log.Println("error finding records")
			return
		}
		var book models.Book
		if err := DB.Where("id = ?",data.BookID).First(&book).Error; err != nil {
			log.Println("error finding records")
			return
		}

		outDate := data.OutDate.Format("01-Jan-2006")
		returnDate:= data.ReturnDate.Format("01-Jan-2006")

		fineMessage := "Subject: Fine Notification - Return Date Exceeded\n\n" +
    "Dear " + user.UserName + ",\n\n" +
    "We hope this email finds you well. We wanted to inform you that the return date for the book \"" + book.BookName + "\" has been exceeded.\n\n" +
    "Book Details:\n" +
    "- Title: " + book.BookName + "\n" +
    "- Received on: " + outDate + "\n" +
    "- Return date: " + returnDate + "\n\n" +
    "A fine of 10 Rs per day has been incurred on your account starting from the due date. To avoid further charges, please return the book as soon as possible.\n\n" +
    "If you have any questions or concerns regarding this fine, please contact our support team.\n\n" +
    "Thank you for your prompt attention to this matter.\n\n" +
    "Best regards,\nThe GoLib Team"

		SendEmail(fineMessage,user.Email)
	}
}


// SendEmail function send the generated message to user email
func SendEmail(message,email string) error {

	SMTPemail := os.Getenv("Email")
	SMTPpass := os.Getenv("Password")
	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}


