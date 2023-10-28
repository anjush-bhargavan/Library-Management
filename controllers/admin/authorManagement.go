package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//GetAuthor handles show Author by id
func GetAuthor(c *gin.Context) {
	id :=c.Param("id")
	var Author models.Author

	if err :=config.DB.First(&Author,id).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{"error" :"Failed to fetch Author"})
		return
	}

	c.JSON(http.StatusOK,Author)
}

// AddAuthors handle admin to add Authors to database
func AddAuthors(c *gin.Context) {
	var Author models.Author

	if err :=c.ShouldBindJSON(&Author); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"error" : "Binding error",
		})
		return
	}
	var existingAuthor models.Author
	if err := config.DB.Where("first_name = ?",Author.FirstName).First(&existingAuthor).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Author already exists"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}

	config.DB.Create(&Author)
	c.JSON(200,gin.H{"message":"Author added succesfully"})
}

//ViewAuthors handles admin to view all Authors in database
func ViewAuthors(c *gin.Context) {

	var Authors []models.Author

	config.DB.Find(&Authors)

	c.JSON(http.StatusOK,Authors)
}

//UpdateAuthor handles admin to update Author details
func UpdateAuthor(c *gin.Context) {
	id :=c.Param("id")
	var Author models.Author

	if err:=config.DB.First(&Author,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{
			"error" : "Author not found",
		})
		return
	}
	if err :=c.ShouldBindJSON(&Author); err!=nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&Author)
	c.JSON(http.StatusOK,Author)
}

//DeleteAuthor handles admin to delete Author by id
func DeleteAuthor(c *gin.Context) {
	id :=c.Param("id")
	var Author models.Author

	if err :=config.DB.First(&Author,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Author not found",
		})
		return
	}
	config.DB.Delete(&Author)
	c.JSON(http.StatusOK,gin.H{"message":"Author deleted succesfully"})
}