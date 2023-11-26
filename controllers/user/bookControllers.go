package controllers

import (
	"net/http"
	"sort"
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
		c.JSON(http.StatusNotFound,gin.H{"status":"Failed",
										"message":"Error fetching book",
										"data":err.Error(),
									})
		return
	}

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Book fetched succesfully",
						"data":book,
					})
}

//ViewBooks handles view all books by pagination
func ViewBooks(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	var books []models.Book

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&books)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":books,
					})
}

//SearchBooks handles users to search books
func SearchBooks(c *gin.Context) {
	query := c.Query("query")
		
	var books []models.Book
	if err := config.DB.Where("book_name ILIKE ?", "%"+query+"%").Find(&books).Error; err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
											"message":"Database error",
											"data":err.Error(),
											})
		return
	}
	
	c.JSON(200, gin.H{	"status":"Success",
						"message":"Search results",
						"data":books,
					})
}


//BookByCategory shows users books by category
func BookByCategory(c *gin.Context) {
	categoryID :=c.Param("id")

	var books models.Book

	if err := config.DB.Where("category_id = ?",categoryID).Find(&books).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
		"message":"Database error",
		"data":err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{	"status":"Success",
						"message":"Search results",
						"data":books,
					})
}


//BookByAuthor shows users books by Author
func BookByAuthor(c *gin.Context) {
	authorID :=c.Param("id")

	var books models.Book

	if err := config.DB.Where("author_id = ?",authorID).Find(&books).Error; err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"status":"Failed",
		"message":"Database error",
		"data":err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{	"status":"Success",
						"message":"Search results",
						"data":books,
					})
}

	//Books type implements sort methods
	type Books []models.Book

	// Len returns the length of the Books slice.
	func (b Books) Len() int {
		return len(b)
	}

	// Swap swaps the elements with indexes i and j.
	func (b Books) Swap(i, j int) {
		b[i], b[j] = b[j], b[i]
	}

	// Less returns true if the book at index i should be ordered before the book at index j.
	// It compares first by OrderCount and then by Rating.
	func (b Books) Less(i, j int) bool {
		// First, compare by OrderCount
		if b[i].OrderCount != b[j].OrderCount {
			return b[i].OrderCount > b[j].OrderCount
		}

		// If OrderCounts are the same, then compare by Rating
		return b[i].Rating > b[j].Rating
	}

//SortByRating function sort the books by ordercount and rating then displays
func SortByRating(c *gin.Context) {
	page,_ :=strconv.Atoi(c.DefaultQuery("page","1"))
	pageSize,_ :=strconv.Atoi(c.DefaultQuery("pageSize","5"))

	books := Books{}

	offset := (page - 1)* pageSize

	config.DB.Order("id").Offset(offset).Limit(pageSize).Find(&books)

	

	sort.Sort(books)

	c.JSON(200,gin.H{	"status":"Success",
						"message":"Books fetched succesfully",
						"data":books,
					})
}
