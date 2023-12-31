package config

import (
	"log"
	"github.com/joho/godotenv"
)

//Loadenv function loads the .env file
func Loadenv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}