package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func GetBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err :=config.DB.First(&book,id).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error" :"Failed to fetch book"})
		return
	}

	c.JSON(http.StatusOK,book)
}


func AddBooks(c *gin.Context) {
	var book models.Book

	if err :=c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"error" : "Binding error",
		})
		return
	}
	var existingBook models.Book
	if err := config.DB.Where("name = ?",book.Book_Name).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Book already exists"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}

	config.DB.Create(&book)
	c.JSON(200,gin.H{"message":"book added succesfully"})
}

func ViewBooks(c *gin.Context) {

	var books []models.Book

	config.DB.Find(&books)

	c.JSON(http.StatusOK,books)
}

func UpdateBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err:=config.DB.First(&book,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{
			"error" : "book not found",
		})
		return
	}
	if err :=c.ShouldBindJSON(&book); err!=nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&book)
	c.JSON(http.StatusOK,book)
}

func DeleteBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err :=config.DB.First(&book,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Book not found",
		})
		return
	}
	config.DB.Delete(&book)
	c.JSON(http.StatusOK,gin.H{"message":"book deleted succesfully"})
}