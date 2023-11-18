package models

// Book model of book details
type Book struct {
	ID         		uint64 			`json:"book_id" gorm:"primaryKey"`
	BookName       	string 			`json:"book_name" gorm:"not null;unique" validate:"required"`
	AuthorID       	uint64 			`json:"author_id" gorm:"not null" validate:"required"`
	NoOfCopies     	uint64 			`json:"no_of_copies" gorm:"not null" validate:"required"`
	Year          	string 			`json:"year" gorm:"not null" validate:"required,len=4"`
	PublicationsID 	uint64 			`json:"publications_id" gorm:"not null" validate:"required"`
	CategoryID     	uint64 			`json:"category_id" gorm:"not null" validate:"required"`
	Description    	string 			`json:"description" gorm:"not null" validate:"required"`
	OrderCount      uint64			`json:"order_count"`
	Rating			uint64			`json:"rating" gorm:"not null;default:5"`
}
