package models

//Publications model of publications
type Publications struct {

	PublicationsID      uint64		`json:"publications_id" gorm:"primaryKey;autoIncrement"`
	PublicationsName    string		`json:"publications_name" gorm:"not null"`
	
}