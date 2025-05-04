package controllers

import (
	"github.com/gin-contrib/sessions"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"weather-app/config"
	"weather-app/models"
	"weather-app/utils"

	"github.com/gin-gonic/gin"
)

// GET /weather
func GetAllWeather(c *gin.Context) {
	session := sessions.Default(c)

	err := session.Get("ERROR")
	deleteSuccess := session.Get("DELETE_SUCCESS")
	errorWeather := session.Get("ERROR_WEATHER")
	successUpdate := session.Get("SUCCESS_UPDATE")

	session.Delete("ERROR")
	session.Delete("DELETE_SUCCESS")
	session.Delete("ERROR_WEATHER")
	session.Delete("SUCCESS_UPDATE")

	menu, _ := utils.GetMenuSubmenu(c)

	var weathers []models.Weather
	config.DB.Preload("Location").Find(&weathers)

	var DataWeathers []map[string]interface{}
	for i, weather := range weathers {
		DataWeathers = append(DataWeathers, map[string]interface{}{
			"Number":        i + 1,
			"ID":            weather.ID,
			"City":          weather.City,
			"Latitude":      weather.Latitude,
			"Longitude":     weather.Longitude,
			"Summary":       weather.Summary,
			"Temperature":   weather.Temperature,
			"WindSpeed":     weather.WindSpeed,
			"WindAngle":     weather.WindAngle,
			"WindDirection": weather.WindDirection,
		})
	}

	c.HTML(http.StatusFound, "show_weather.html", gin.H{
		"title":         "Weathers",
		"menu":          menu,
		"data":          DataWeathers,
		"user":          session.Get("USERNAME"),
		"errorWeather":  errorWeather,
		"err":           err,
		"successUpdate": successUpdate,
		"deleteSuccess": deleteSuccess,
	})
}

// GET /weather/:id
func GetWeatherByID(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	menu, _ := utils.GetMenuSubmenu(c)
	var weather models.Weather

	if err := config.DB.First(&weather, "id = ?", id).Error; err != nil {
		session.Set("ERROR_WEATHER", "Weather not found")
		c.Redirect(http.StatusFound, "home/weather/")
		return
	}

	c.HTML(http.StatusOK, "detail_weather.html", gin.H{
		"title": "Detail Weather",
		"menu":  menu,
		"data":  weather,
	})
}

// GET /weather/edit/:id
func EditWeatherByID(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	menu, _ := utils.GetMenuSubmenu(c)
	var weather models.Weather

	if err := config.DB.First(&weather, "id = ?", id).Error; err != nil {
		session.Set("ERROR_WEATHER", "Weather not found")
		c.Redirect(http.StatusFound, "home/weather/")
		return
	}

	c.HTML(http.StatusOK, "edit_weather.html", gin.H{
		"title":     "Edit Weather",
		"menu":      menu,
		"data":      weather,
		"csrfToken": csrf.GetToken(c),
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
	session := sessions.Default(c)
	id := c.Param("id")
	var weather models.Weather

	if err := config.DB.First(&weather, id).Error; err != nil {
		session.Set("ERROR", "Weather not found")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	var input models.Weather
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	// Update manual field
	weather.Summary = input.Summary
	weather.Temperature = input.Temperature
	weather.WindSpeed = input.WindSpeed
	weather.WindAngle = input.WindAngle
	weather.WindDirection = input.WindDirection

	config.DB.Save(&weather)

	session.Set("SUCCESS_UPDATE", "Weather successfully updated")
	session.Save()
	c.Redirect(http.StatusFound, "/home/weather")
}

// DELETE /weather/del:id
func DeleteWeather(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)

	if err := config.DB.Delete(&models.Weather{}, id).Error; err != nil {
		session.Set("DELETE_ERROR", "Delete failed")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	session.Set("DELETE_SUCCESS", "Weather successfully deleted")
	session.Save()
	c.Redirect(http.StatusFound, "/home/weather")
}
