package models

//FeedBack model holds the detail of user feedback
type FeedBack struct{
	ID			uint64			`json:"feedback_id" gorm:"primaryKey"`
	UserID		uint64			`json:"user_id" gorm:"not null"`
	Subject		string			`json:"subject" gorm:"not null"`
	Content		string			`json:"content"`
} 