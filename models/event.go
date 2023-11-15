package models


//Event model holds the details of events organized by library
type Event struct {
	ID			uint64			`json:"event_id" gorm:"primaryKey"`
	Name		string			`json:"name" gorm:"not null"`
	Date		string			`json:"date" gorm:"not null"`
	Time		string			`json:"time"`
	Venue		string			`json:"venue"`
	Description	string			`json:"description"`	
}