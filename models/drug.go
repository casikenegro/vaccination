package models

import (
	"time"

	"gorm.io/gorm"
)

type Drug struct {
	gorm.Model
	Id           uint          `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name"`
	Approved     bool          `json:"approved"`
	MinDose      int64         `json:"min_dose"`
	MaxDose      int64         `json:"max_dose"`
	AvaliableAt  time.Time     `json:"avaliable_at"`
	Vaccinations []Vaccination `json:"vaccinations" gorm:"foreignKey:DrugId"`
}
