package models

//Author model of author details
type Author struct {
	ID      	  uint64		`json:"author_id" gorm:"primaryKey"`
	FirstName	  string		`json:"first_name" gorm:"not null"`
	LastName      string		`json:"last_name" gorm:"not null"`	
	Language      string		`json:"language" gorm:"not null"`
	Remarks       string		`json:"remarks"`

} 