package models

import "gorm.io/gorm"

type Weather struct {
	gorm.Model
	LocationID    uint    `json:"location_id"`
	Summary       string  `json:"weather_summary"`
	Temp          float64 `json:"temperature"`
	WindSpeed     float64 `json:"wind_speed"`
	WindAngle     float64 `json:"wind_angle"`
	WindDirection string  `json:"wind_direction"`
}
