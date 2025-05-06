package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey" json:"id"`
	City      string  `gorm:"unique" json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Weather   Weather `gorm:"foreignKey:LocationID"`
}
