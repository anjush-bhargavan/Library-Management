package models

//Cart model contains the id of books in cart
type Cart struct {
	BookID uint64 `json:"book_id" gorm:"foreignKey:BookID;not null;unique"`
}