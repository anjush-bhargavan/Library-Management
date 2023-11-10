package models

//Publications model of publications
type Publications struct {

	ID      			uint64		`json:"publications_id" gorm:"primaryKey"`
	PublicationsName    string		`json:"publications_name" gorm:"not null"`

}