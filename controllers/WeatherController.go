package controllers

import (
	"net/http"
	"weather-app/config"
	"weather-app/models"

	"github.com/gin-gonic/gin"
)

// GET /weather
func GetAllWeather(c *gin.Context) {
	var weathers []models.Weather
	config.DB.Preload("Location").Find(&weathers)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "All weather data",
		"data":    weathers,
	})
}

// GET /weather/:id
func GetWeatherByID(c *gin.Context) {
	id := c.Param("id")
	var weather models.Weather

	if err := config.DB.First(&weather, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Weather not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Weather found",
		"data":    weather,
	})
}

// POST /weather
func CreateWeather(c *gin.Context) {
	var input models.Weather

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save weather"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Weather created",
		"data":    input,
	})
}

// PUT /weather/:id
func UpdateWeather(c *gin.Context) {
	id := c.Param("id")
	var weather models.Weather

	if err := config.DB.First(&weather, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Weather not found"})
		return
	}

	var input models.Weather
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Update manual field
	weather.Summary = input.Summary
	weather.Temperature = input.Temperature
	weather.WindSpeed = input.WindSpeed
	weather.WindAngle = input.WindAngle
	weather.WindDirection = input.WindDirection

	config.DB.Save(&weather)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Weather updated",
		"data":    weather,
	})
}

// DELETE /weather/:id
func DeleteWeather(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Weather{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Weather deleted"})
}
