package models


//User model of user details
type User struct{
	UserID       uint64         `json:"user_id" gorm:"primaryKey"`
	FirstName    string		    `json:"first_name" gorm:"not null" validate:"required"` 		
	LastName     string		    `json:"last_name" gorm:"not null" validate:"required"`
	UserName     string			`json:"user_name" gorm:"not null"`
	DoB          string	        `json:"date_of_birth" gorm:"not null" validate:"required"`
	Gender       string			`json:"gender" gorm:"not null;check gender IN('M','F','other')" validate:"required"`
	Email        string			`json:"email" gorm:"not null;unique" validate:"required,email"`
	Phone        string			`json:"phone" gorm:"not null;unique" validate:"required,len=10"`
	Role 		 string			`json:"role" gorm:"not null;default:'user'"`
	Address      string			`json:"address" gorm:"not null" validate:"required"`
	Password     string			`json:"password" gorm:"not null" validate:"required"`
}



