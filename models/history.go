package models

import "time"

//History model the data that rented
type History struct {
	ID        		 uint64			`json:"history_id" gorm:"primaryKey"`
	UserID           uint64			`json:"user_id" gorm:"not null"`
	BookID           uint64			`json:"book_id" gorm:"not null"`
	Status	 		 string			`json:"status" gorm:"not null;default:'pending'"`
	OrderedOn		 time.Time		`json:"ordered_on" gorm:"not null"`
	RentedOn         time.Time		`json:"rented_on"`
	ReturnedOn       time.Time		`json:"returned_on"`
	Remarks          string			`json:"remarks"`
}