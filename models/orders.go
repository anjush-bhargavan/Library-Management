package models

import "time"

//Orders holds the details of orders
type Orders struct{
	ID        		 uint64			`json:"order_id" gorm:"primaryKey"`
	UserID           uint64			`json:"user_id" gorm:"not null"`
	BookID           uint64			`json:"book_id" gorm:"not null"`
	Type			 string			`json:"type"`
	Status	 		 string			`json:"status" gorm:"not null;default:'pending'"`
	OrderedOn		 time.Time		`json:"ordered_on" gorm:"not null"`
	Remarks          string			`json:"remarks"`
}