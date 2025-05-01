package controllers

import (
	"net/http"
	"weather-app/config"
	"weather-app/models"

	"github.com/gin-gonic/gin"
)

func SaveLocation(c *gin.Context) {
	var input struct {
		City      string  `json:"city" binding:"required"`
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Cek duplicate berdasarkan city + latitude + longitude
	var existing models.Location
	if err := config.DB.Where("city = ? AND latitude = ? AND longitude = ?",
		input.City, input.Latitude, input.Longitude).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "Location already exists",
		})
		return
	}

	// Simpan data
	location := models.Location{
		City:      input.City,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Weather: models.Weather{
			Summary:       "",
			Temp:          0,
			WindSpeed:     0,
			WindAngle:     0,
			WindDirection: "",
		},
	}

	if err := config.DB.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save location"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Success create weather data",
		"data": gin.H{
			"ID":              location.ID,
			"CreatedAt":       location.CreatedAt,
			"UpdatedAt":       location.UpdatedAt,
			"DeletedAt":       location.DeletedAt,
			"city":            location.City,
			"latitude":        location.Latitude,
			"longitude":       location.Longitude,
			"weather_summary": location.Weather.Summary,
			"temperature":     location.Weather.Temp,
			"wind_speed":      location.Weather.WindSpeed,
			"wind_angle":      location.Weather.WindAngle,
			"wind_direction":  location.Weather.WindDirection,
		},
	})
}
