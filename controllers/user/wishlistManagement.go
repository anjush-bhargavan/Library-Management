package controllers

import (
	"net/http"
	"strconv"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//AddToWishlist function add books to Wishlist
func AddToWishlist(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)

	stringID :=c.Param("id")
	var Wishlist models.Wishlist
		id, err := strconv.ParseUint(stringID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest,gin.H{"error":"error parsing string"})
        return
    }
	var book models.Book
	if err := config.DB.Where("id = ?",id).First(&book).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"book not found in records"})
		return
	}
	Wishlist.BookID=id
	var existingWishlist models.Wishlist
	if err := config.DB.Where("book_id = ? AND user_id = ?",Wishlist.BookID,userID).First(&existingWishlist).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Book already added to Wishlist"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}
	
	Wishlist.UserID=userID
	config.DB.Create(&Wishlist)
	c.JSON(http.StatusOK,gin.H{"message":"book added to Wishlist"})

}


//ShowWishlist function lists the Wishlist items
func ShowWishlist(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	var books []models.Book
	var bookIds []uint64
	if err :=config.DB.Model(&models.Wishlist{}).Where("user_id = ?",userID).Pluck("BookID",&bookIds).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"database error"})
		return
	}
	if len(bookIds) == 0 {
		c.JSON(http.StatusNoContent,gin.H{"error":"No items in Wishlist"})
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


//DeleteWishlist function deletes the Wishlist items
func DeleteWishlist(c *gin.Context) {
	userIDContext,_ :=c.Get("user_id")
	userID:=userIDContext.(uint64)
	id :=c.Param("id")
	var book models.Wishlist

	if err := config.DB.Where("book_id = ? AND user_id = ?",id,userID).Delete(&book).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"error":"error deleting from database"})
	}
	c.JSON(http.StatusOK,gin.H{"message":"book removed from Wishlist"})

}