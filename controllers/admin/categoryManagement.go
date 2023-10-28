package controllers

import (
	"fmt"
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCategory handles show Category by id
func GetCategory(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	var category models.Category

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch Category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// AddCategorys handle admin to add Categorys to database
func AddCategorys(c *gin.Context) {
	var Category models.Category

	if err := c.ShouldBindJSON(&Category); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Binding error",
		})
		return
	}
	var existingCategory models.Category
	if err := config.DB.Where("category_name = ?", Category.CategoryName).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	config.DB.Create(&Category)
	c.JSON(200, gin.H{"message": "Category added succesfully"})
}

// ViewCategorys handles admin to view all Categorys in database
func ViewCategorys(c *gin.Context) {

	var Categorys []models.Category

	config.DB.Find(&Categorys)

	c.JSON(http.StatusOK, Categorys)
}

// UpdateCategory handles admin to update Category details
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var Category models.Category

	if err := config.DB.First(&Category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}
	if err := c.ShouldBindJSON(&Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&Category)
	c.JSON(http.StatusOK, Category)
}

// DeleteCategory handles admin to delete Category by id
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var Category models.Category

	if err := config.DB.First(&Category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}
	config.DB.Delete(&Category)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted succesfully"})
}
