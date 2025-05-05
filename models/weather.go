package models

import "gorm.io/gorm"

type Weather struct {
	gorm.Model
	LocationID    uint      `gorm:"uniqueIndex"`
	City          string    `json:"city"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Summary       string    `json:"weather_summary"`
	Temperature   float64   `json:"temperature"`
	WindSpeed     float64   `json:"wind_speed"`
	WindAngle     float64   `json:"wind_angle"`
	WindDirection string    `json:"wind_direction"`
	Location      *Location `gorm:"foreignKey:LocationID"`
}
