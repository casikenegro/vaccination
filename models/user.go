package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Id       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"not null, unique"`
	Password string `json:"password" gorm:"not null"`
}
