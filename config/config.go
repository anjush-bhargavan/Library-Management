package config

import (
	"log"
	"os"

	"github.com/anjush-bhargavan/library-management/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectDB(){
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DB_Config")
	var err error

	DB,err = gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Author{},
		&models.Category{},
		&models.Publications{},
	)
}