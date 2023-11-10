package models

//Wishlist model contains the id of books in cart
type Wishlist struct {
	ID		uint64	`json:"wishlist_id" gorm:"primaryKey"`
	UserID  uint64 	`json:"user_id" gorm:"not null"`  
	BookID 	uint64	`json:"book_id" gorm:"not null"`
}