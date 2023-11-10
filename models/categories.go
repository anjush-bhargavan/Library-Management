package models

//Category model of category
type Category struct {
    ID              uint64 `json:"category_id" gorm:"primaryKey"`
    CategoryName    string `json:"category_name" gorm:"not null"`

}
