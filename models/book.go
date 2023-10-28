package models

// Book model of book details
type Book struct {
	BookID         uint64 `json:"book_id" gorm:"primaryKey;autoIncrement"`
	BookName       string `json:"book_name" gorm:"not null;unique"`
	AuthorID       uint64 `json:"author_id" gorm:"foreignKey:author_id;not null"`
	NoOfCopies     uint64 `json:"no_of_copies" gorm:"not null"`
	Year           uint64 `json:"year" gorm:"not null"`
	PublicationsID uint64 `json:"publications_id" gorm:"foreignKey:publications_id;not null"`
	CategoryID     uint64 `json:"category_id" gorm:"not null"`
	Description    string `json:"description" gorm:"not null"`
	ReviewID       uint64 `json:"review_id" gorm:"not null"`
}
