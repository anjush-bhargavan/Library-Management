package config

import (
	"fmt"
	"log"
	"time"

	"github.com/anjush-bhargavan/library-management/models"
	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New()

	c.AddFunc("@daily", CheckMembership)

	c.Start()
}

func CheckMembership() {
	fmt.Println("here i am")
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
