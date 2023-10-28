package models


//User model of user details
type User struct{
	UserID       uint64         `json:"user_id" gorm:"primaryKey;autoIncrement"`
	FirstName    string		    `json:"first_name" gorm:"not null"` 		
	LastName     string		    `json:"last_name" gorm:"not null"`
	Gender       string			`json:"gender" gorm:"not null;check gender IN('M','F','other')"`
	Email        string			`json:"email" gorm:"not null;unique"`
	Phone        uint64			`json:"phone" gorm:"not null;unique;check:phone_length=10"`
	Role 		 string			`json:"role" gorm:"not null;default:'user'"`
	Address      string			`json:"address" gorm:"not null"`
	Password     string			`json:"password" gorm:"not null"`
}



