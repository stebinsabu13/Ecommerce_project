package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstname" gorm:"not null"`
	LastName  string `json:"lastname"  gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	MobileNum string `json:"mobilenum" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" gorm:"not null"`
	Role      string `json:"role" gorm:"default:'USER'"`
}
