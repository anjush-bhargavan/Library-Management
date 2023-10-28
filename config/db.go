package config

import (
	"os"
	"github.com/anjush-bhargavan/library-management/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB variable globally declared
var DB *gorm.DB

//ConnectDB to connect database
func ConnectDB(){
	dsn := os.Getenv("DB_Config")
	var err error
	DB,err = gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	DB.AutoMigrate(
		&models.Category{},
		&models.Author{},
		&models.Publications{},
		&models.Cart{},
		&models.User{},
		&models.Book{},
	)
}