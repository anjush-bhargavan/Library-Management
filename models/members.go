package models

import "time"

//Membership model holds the record of subcription
type Membership struct {
	ID			    		uint64          `json:"membership_id" gorm:"PrimaryKey"`
	UserID          		uint64			`json:"user_id" gorm:"not null;unique;foreignKey:UserID"`
	RazorpaySubscriptionID	string			`json:"subscription_id" gorm:"not null;unique"`
	Plan 					string			`json:"plan" gorm:"not null"`
	StartedOn       		time.Time		`json:"started_on" gorm:"not null"`
	ExpiresAt       		time.Time		`json:"expires_at" gorm:"not null"`

}