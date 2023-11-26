package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anjush-bhargavan/library-management/config"
	"github.com/anjush-bhargavan/library-management/models"
	"github.com/gin-gonic/gin"
)

//ViewMemberships function list out all memberships
func ViewMemberships(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var users []models.Membership

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&users)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":users,
					})
}

//RemoveMembership handles admin to edit members
func RemoveMembership(c *gin.Context) {
	id :=c.Param("id")
	userID,_ :=strconv.Atoi(id)
	var member models.Membership

	if err := config.DB.Where("user_id = ?",userID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	member.IsActive=false
	if err := config.DB.Save(&member).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Membership details",
						"data":member,
						})

}

//GetMembership functions gets the details of member
func GetMembership(c *gin.Context) {
	id :=c.Param("id")
	userID,_ :=strconv.Atoi(id)
	var member models.Membership

	if err := config.DB.Where("user_id = ?",userID).First(&member).Error; err != nil {
		c.JSON(http.StatusNotFound,gin.H{	"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
										})
		return
	}
	var result struct{
		ID			uint64
		UserName	string
		Phone		string
		Address		string
		Email		string
		Plan		string
		IsActive	bool
		StartedOn	time.Time
		ExpiresAt	time.Time
	}
	query := `
    SELECT 
        memberships.id,
        users.user_name,
        users.phone,
        users.address,
		users.email,
        memberships.plan,
        memberships.is_active,
        memberships.started_on,
		memberships.expires_at
    FROM 
        memberships
    JOIN 
        users ON memberships.user_id = users.user_id
    WHERE 
        membership.id = ?`

	if err := config.DB.Raw(query, member.UserID).Scan(&result).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{	"status":"Failed",
												"message":"Database error while joining",
												"data":err.Error(),
											})
		return
	}
	c.JSON(200,gin.H{	"status":"Success",
						"message":"Membership details",
						"data":result,
						})
}