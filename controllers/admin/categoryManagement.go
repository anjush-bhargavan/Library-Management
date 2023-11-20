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
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Failed to fetch Category",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Category fetched succesfully",
						"data":category,
					})
}

// AddCategorys handle admin to add Categorys to database
func AddCategorys(c *gin.Context) {
	var Category models.Category

	if err := c.ShouldBindJSON(&Category); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{	"status":"Failed",
												"message":"Binding error",
												"data":err.Error(),
											})
		return
	}

	if err := validate.Struct(Category); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	var existingCategory models.Category
	if err := config.DB.Where("category_name = ?", Category.CategoryName).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{	"status":"Failed",
											"message":"Category already exists",
											"data":existingCategory,
										})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{	"status":"Failed",
														"message":"Database error",
														"data":err.Error(),
													})
		return
	}

	config.DB.Create(&Category)
	c.JSON(200, gin.H{	"status":"Success",
						"message":"Category added succesfully",
						"data":Category,
					})
}

// ViewCategorys handles admin to view all Categorys in database
func ViewCategorys(c *gin.Context) {

	var Categories []models.Category

	config.DB.Find(&Categories)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Categories fetched succesfully",
						"data":Categories,
					})
}

// UpdateCategory handles admin to update Category details
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var Category models.Category

	if err := config.DB.First(&Category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Category not found",
											"data":err.Error(),
										})
		return
	}
	if err := c.ShouldBindJSON(&Category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{	"status":"Failed",
												"message":"Binding error",
												"data":err.Error(),
											})
		return
	}

	if err := validate.Struct(Category); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	config.DB.Save(&Category)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"category updated succesfully",
						"data":Category,
					})
}

// DeleteCategory handles admin to delete Category by id
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var Category models.Category

	if err := config.DB.First(&Category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Category not found",
											"data":err.Error(),
										})
		return
	}
	config.DB.Delete(&Category)
	c.JSON(http.StatusOK, gin.H{	"status":"Success",
									"message":"Category deleted succesfully",
									"data":Category,
								})
}
