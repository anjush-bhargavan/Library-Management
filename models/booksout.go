package models

import "time"

//BooksOut model contains details of book currently out
type BooksOut struct {
	ID    		uint64			`json:"booksout_id" gorm:"primaryKey"`
	UserID       uint64 		`json:"user_id" gorm:"not null"`
	BookID       uint64			`json:"book_id" gorm:"not null"`
	OutDate      time.Time		`json:"out_date" gorm:"not null"`
	ReturnDate	 time.Time		`json:"return_date" gorm:"not null"`

}