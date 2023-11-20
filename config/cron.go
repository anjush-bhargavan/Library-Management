package config

import (
	"log"
	"time"

	"github.com/anjush-bhargavan/library-management/models"
	"github.com/robfig/cron/v3"
)


//InitCron function initializes the cron job and run the CheckMembership in each day
func InitCron() {
	c := cron.New()

	c.AddFunc("@daily", CheckMembership)

	c.Start()
}

//CheckMembership will check users membership is expired or not in each day
func CheckMembership() {
	var members []models.Membership
	now := time.Now()
	if err := DB.Where("expires_at < ?", now).Find(&members).Error; err != nil {
		log.Println("error finding records")
		return
	}
	for _, user := range members {
		user.IsActive = false
		DB.Save(&user)
	}

}
