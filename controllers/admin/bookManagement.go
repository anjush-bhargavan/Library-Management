package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


//GetBook handles show book by id
func GetBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err :=config.DB.First(&book,id).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error" :"Failed to fetch book"})
		return
	}

	c.JSON(http.StatusOK,book)
}

// AddBooks handle admin to add books pto database
func AddBooks(c *gin.Context) {
	var book models.Book

	if err :=c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"error" : "Binding error",
		})
		c.Abort()
		return
	}

	if err := validate.Struct(book); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Please fill all fields"+err.Error()})
		c.Abort()
		return
	}

	var existingBook models.Book
	if err := config.DB.Where("book_name = ?",book.BookName).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Book already exists"})
		c.Abort()
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		c.Abort()
		return
	}

	if err:=config.DB.Create(&book).Error; err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Error adding to database"})
		return
	}
	c.JSON(200,gin.H{"message":"book added succesfully"})
}

//ViewBooks handles admin to view all books in database
func ViewBooks(c *gin.Context) {

	var books []models.Book

	config.DB.Find(&books)

	c.JSON(http.StatusOK,books)
}

//UpdateBook handles admin to update book details
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

	if err := validate.Struct(book); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Please fill all fields"})
		return
	}

	config.DB.Save(&book)
	c.JSON(http.StatusOK,book)
}

//DeleteBook handles admin to delete book by id
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