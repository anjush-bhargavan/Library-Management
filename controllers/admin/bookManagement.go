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
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Failed to fetch book",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"book fetched succesfully",
						"data":book,
					})
}

// AddBooks handle admin to add books pto database
func AddBooks(c *gin.Context) {
	var book models.Book

	if err :=c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
											})
		c.Abort()
		return
	}

	if err := validate.Struct(book); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		c.Abort()
		return
	}

	var existingBook models.Book
	if err := config.DB.Where("book_name = ?",book.BookName).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Book already exists",
											"data":existingBook,
										})
		c.Abort()
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{	"status":"Failed",
														"message":"Database error",
														"data":err.Error(),
													})
		c.Abort()
		return
	}

	if err:=config.DB.Create(&book).Error; err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Error adding to database",
											"data":err.Error(),
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Book added succesfully",
						"data":book,
					})
}

//ViewBooks handles admin to view all books in database
func ViewBooks(c *gin.Context) {

	var books []models.Book

	config.DB.Find(&books)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":books,
					})
}

//UpdateBook handles admin to update book details
func UpdateBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err:=config.DB.First(&book,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Book not found",
											"data":err.Error(),
										})
		return
	}
	if err :=c.ShouldBindJSON(&book); err!=nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}

	if err := validate.Struct(book); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	config.DB.Save(&book)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Book updated succesfully",
						"data":book,
					})
}

//DeleteBook handles admin to delete book by id
func DeleteBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err :=config.DB.First(&book,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Book not found",
											"data":err.Error(),
										})
		return
	}
	config.DB.Delete(&book)
	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Book deleted succesfully",
								"data":nil,
							})
}