package models

import (
	"time"

	"gorm.io/gorm"
)

type Vaccination struct {
	gorm.Model

	Id     uint      `json:"id" gorm:"primaryKey"`
	Name   string    `json:"name"`
	Dose   int64     `json:"dose"`
	Date   time.Time `json:"date"`
	DrugId uint      `json:"drug_id"`
}
