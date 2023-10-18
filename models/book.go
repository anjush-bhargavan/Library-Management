package models

type Book struct {
	Book_Id          uint64     `json:"book_id" gorm:"primaryKey;autoIncrement"`
	Book_Name        string 	`json:"book_name" gorm:"not null;unique"`
	Author_Id        uint64		`json:"author_id" gorm:"not null"`
	No_Of_Copies     uint64		`json:"no_of_copies" gorm:"not null"`
	Year             uint64		`json:"year" gorm:"not null"`
	Publications_Id  uint64		`json:"publications_id" gorm:"not null"`
	Category_Id      uint64		`json:"category_id" gorm:"not null"`
	Description      string		`json:"description" gorm:"not null"`
	Review_Id        uint64		`json:"review_id" gorm:"not null"`
}