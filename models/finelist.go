package models


//FineList model contains fine of users
type FineList struct {
	ID	    uint64		`json:"fine_id" gorm:"primaryKey"`
	UserID  uint64		`json:"user_id" gorm:"not null;unique"`
	Fine    uint64		`json:"fine"`
}
