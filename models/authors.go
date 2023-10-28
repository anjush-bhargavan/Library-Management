package models

//Author model of author details
type Author struct {
	AuthorID      uint64		`json:"author_id" gorm:"primaryKey;autoIncrement"`
	FirstName	  string		`json:"first_name" gorm:"not null"`
	LastName      string		`json:"last_name" gorm:"not null"`	
	Language      string		`json:"language" gorm:"not null"`
	Remarks       string		`json:"remarks"`

} 