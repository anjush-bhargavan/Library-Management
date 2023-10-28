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
		c.JSON(http.StatusNotFound,gin.H{"error" :"Failed to fetch Publications"})
		return
	}

	c.JSON(http.StatusOK,Publications)
}

// AddPublication handle admin to add Publications pto database
func AddPublication(c *gin.Context) {
	var Publications models.Publications

	if err :=c.ShouldBindJSON(&Publications); err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"error" : "Binding error",
		})
		return
	}
	var existingPublications models.Publications
	if err := config.DB.Where("publications_name = ?",Publications.PublicationsName).First(&existingPublications).Error; err == nil {
		c.JSON(http.StatusConflict,gin.H{"error":"Publications already exists"})
		return
	}else if err !=  gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Database error"})
		return
	}

	config.DB.Create(&Publications)
	c.JSON(200,gin.H{"message":"Publications added succesfully"})
}

//ViewPublications handles admin to view all Publications in database
func ViewPublications(c *gin.Context) {

	var Publications []models.Publications

	config.DB.Find(&Publications)

	c.JSON(http.StatusOK,Publications)
}

//UpdatePublication handles admin to update Publications details
func UpdatePublication(c *gin.Context) {
	id :=c.Param("id")
	var Publications models.Publications

	if err:=config.DB.First(&Publications,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{
			"error" : "Publications not found",
		})
		return
	}
	if err :=c.ShouldBindJSON(&Publications); err!=nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&Publications)
	c.JSON(http.StatusOK,Publications)
}

//DeletePublication handles admin to delete Publications by id
func DeletePublication(c *gin.Context) {
	id :=c.Param("id")
	var Publications models.Publications

	if err :=config.DB.First(&Publications,id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":"Publications not found",
		})
		return
	}
	config.DB.Delete(&Publications)
	c.JSON(http.StatusOK,gin.H{"message":"Publications deleted succesfully"})
}