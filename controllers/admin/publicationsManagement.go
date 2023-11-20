package controllers

import (
	"net/http"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


//GetPublication handles show Publications by id
func GetPublication(c *gin.Context) {
	id :=c.Param("id")
	var Publications models.Publications

	if err :=config.DB.First(&Publications,id).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Failed to fetch Publications",
											"data":err.Error(),
										})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Publication fetched succesfully",
						"data":Publications,
					})
}

// AddPublication handle admin to add Publications pto database
func AddPublication(c *gin.Context) {
	var Publications models.Publications

	if err :=c.ShouldBindJSON(&Publications); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}

	if err := validate.Struct(Publications); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}

	var existingPublications models.Publications
	if err := config.DB.Where("publications_name = ?",Publications.PublicationsName).First(&existingPublications).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{	"status":"Failed",
											"message":"Publications already exists",
											"data":existingPublications,
										})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{	"status":"Failed",
														"message":"Database error",
														"data":err.Error(),
													})
		return
	}

	config.DB.Create(&Publications)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Publications added succesfully",
						"data":Publications,
					})
}

//ViewPublications handles admin to view all Publications in database
func ViewPublications(c *gin.Context) {

	var Publications []models.Publications

	config.DB.Find(&Publications)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Publications fetched succesfully",
						"data":Publications,
					})
}

//UpdatePublication handles admin to update Publications details
func UpdatePublication(c *gin.Context) {
	id :=c.Param("id")
	var Publications models.Publications

	if err:=config.DB.First(&Publications,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Publications not found",
											"data":err.Error(),
										})
		return
	}
	if err :=c.ShouldBindJSON(&Publications); err!=nil {
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Binding error",
											"data":err.Error(),
										})
		return
	}

	if err := validate.Struct(Publications); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{	"status":"Failed",
											"message":"Please fill all fields",
											"data":err.Error(),
										})
		return
	}
	config.DB.Save(&Publications)
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Publications updated succesfully",
						"data":Publications,
					})
}

//DeletePublication handles admin to delete Publications by id
func DeletePublication(c *gin.Context) {
	id :=c.Param("id")
	var Publications models.Publications

	if err :=config.DB.First(&Publications,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{	"status":"Failed",
											"message":"Publications not found",
											"data":err.Error(),
										})
		return
	}
	config.DB.Delete(&Publications)
	c.JSON(http.StatusOK,gin.H{	"status":"Success",
								"message":"Publications deleted succesfully",
								"data":nil,
							})
}