package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate =validator.New()

//GetAuthor handles show Author by id
func GetAuthor(c *gin.Context) {
	id :=c.Param("id")
	var Author models.Author

	if err :=config.DB.First(&Author,id).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Failed to fetch Author",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Author fetched succesfully",
						"data":Author,
					})
}

// AddAuthors handle admin to add Authors to database
func AddAuthors(c *gin.Context) {
	var Author models.Author

	if err :=c.ShouldBindJSON(&Author); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}

	if err := validate.Struct(Author); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	var existingAuthor models.Author
	if err := config.DB.Where("first_name = ?",Author.FirstName).First(&existingAuthor).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Author already exists",
											"data":existingAuthor,
										})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{	"status":"Failed",
														"message":"Database error",
														"data":err.Error(),
													})
		return
	}

	config.DB.Create(&Author)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Author added succesfully",
						"data":Author,
					})
}

//ViewAuthors handles admin to view all Authors in database
func ViewAuthors(c *gin.Context) {

	var Authors []models.Author

	config.DB.Find(&Authors)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Authors fetched succesfully",
						"data":Authors,
					})
}

//UpdateAuthor handles admin to update Author details
func UpdateAuthor(c *gin.Context) {
	id :=c.Param("id")
	var Author models.Author

	if err:=config.DB.First(&Author,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Author not found",
											"data":err.Error(),
										})
		return
	}
	if err :=c.ShouldBindJSON(&Author); err!=nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}

	if err := validate.Struct(Author); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	config.DB.Save(&Author)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Authors updated succesfully",
						"data":Author,
					})
}

//DeleteAuthor handles admin to delete Author by id
func DeleteAuthor(c *gin.Context) {
	id :=c.Param("id")
	var Author models.Author

	if err :=config.DB.First(&Author,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Author not found",
											"data":err.Error(),
										})
		return
	}
	config.DB.Delete(&Author)
	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Author deleted succesfully",
								"data":nil,
							})
}