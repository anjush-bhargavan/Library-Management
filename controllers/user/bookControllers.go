package controllers

import (
	"net/http"
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//GetBook handles get book by id
func GetBook(c *gin.Context) {
	id :=c.Param("id")
	var book models.Book

	if err :=config.DB.First(&book,id).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error" :"Failed to fetch book"})
		return
	}

	c.JSON(http.StatusOK,book)
}

//ViewBooks handles view all books by pagination
func ViewBooks(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var books []models.Book

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&books)

	c.JSON(http.StatusOK,books)
}