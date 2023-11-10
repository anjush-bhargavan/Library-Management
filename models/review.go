package models

//Review model holds the reviews of books by user
type Review struct {
	ID    		uint64		`json:"review_id" gorm:"primaryKey"`
	UserID      uint64		`json:"user_id" gorm:"not null"`
	BookID      uint64		`json:"book_id" gorm:"not null"`
	Rating      uint64		`json:"rating" gorm:"not null"`
	Review      string		`json:"review"`
}
