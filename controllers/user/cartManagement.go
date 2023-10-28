package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//AddToCart function add books to cart
func AddToCart(c *gin.Context) {
	stringID :=c.Param("id")
	var book models.Cart
		id, err := strconv.ParseUint(stringID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest,gin.H{"error":"error parsing string"})
        return
    }
	book.BookID=id
	var existingBook models.Cart
	if err := config.DB.Where("book_id = ?",book.BookID).First(&existingBook).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Book already added to cart"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}
	config.DB.Create(&book)
	c.JSON(http.StatusOK,gin.H{"message":"book added to cart"})

}


//ShowCart function lists the cart items
func ShowCart(c *gin.Context) {
	var books []models.Book
	var bookIds []uint64
	config.DB.Model(&models.Cart{}).Pluck("BookID",&bookIds)
	fmt.Println(bookIds)
	if len(bookIds) == 0 {
		c.JSON(http.StatusNoContent,gin.H{"error":"No items in cart"})
		return
	}
	for _,id := range bookIds{
		var book models.Book

		if err :=config.DB.First(&book,id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":"book not found",
			})
			return
		}
		books=append(books, book)
	}

	c.JSON(http.StatusOK,books)
}


//DeleteCart function deletes the cart items
func DeleteCart(c *gin.Context) {
	id :=c.Param("id")
	var book models.Cart

	if err := config.DB.Where("book_id = ?",id).Delete(&book).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"error deleting from database"})
	}
	c.JSON(http.StatusOK,gin.H{"message":"book removed from cart"})

}