package models

//Category model of category
type Category struct {
    CategoryID   uint64 `json:"category_id" gorm:"primaryKey;autoIncrement"`
    CategoryName string `json:"category_name" gorm:"not null"`
}
