package controllers

import (
	"github.com/gin-contrib/sessions"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"weather-app/config"
	"weather-app/middlewares"
	"weather-app/models"
	"weather-app/utils"

	"github.com/gin-gonic/gin"
)

// GET /weather
func GetAllWeather(c *gin.Context) {
	session := sessions.Default(c)
	userSession := middlewares.GetSessionUser(c)
	session.Set("USERNAME_SESSION", userSession.Username)
	session.Save()

	err := utils.FlashMessage(c, "ERROR")
	successDelete := utils.FlashMessage(c, "DELETE_SUCCESS")
	errorWeather := utils.FlashMessage(c, "ERROR_WEATHER")
	successUpdate := utils.FlashMessage(c, "SUCCESS_UPDATE")
	errorUpdate := utils.FlashMessage(c, "ERROR_UPDATE")
	deleteError := utils.FlashMessage(c, "DELETE_ERROR")

	menu, _ := utils.GetMenuSubmenu(c)

	var weathers []models.Weather
	config.DB.Preload("Location").Find(&weathers)

	var DataWeathers []map[string]interface{}
	for i, weather := range weathers {
		DataWeathers = append(DataWeathers, map[string]interface{}{
			"Number":        i + 1,
			"ID":            utils.EncodeID(int(weather.ID)),
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
		"title":         "Weather",
		"menu":          menu,
		"data":          DataWeathers,
		"user":          userSession,
		"errorWeather":  errorWeather,
		"err":           err,
		"successUpdate": successUpdate,
		"errorUpdate":   errorUpdate,
		"successDelete": successDelete,
		"deleteError":   deleteError,
	})
}

// GET /weather/:id
func GetWeatherByID(c *gin.Context) {
	id := c.Param("id")
	DecodedID, _ := utils.DecodeID(id)
	session := sessions.Default(c)
	userSession := middlewares.GetSessionUser(c)
	menu, _ := utils.GetMenuSubmenu(c)
	var weather models.Weather

	session.Set("USERNAME_SESSION", userSession.Username)
	session.Save()

	if err := config.DB.First(&weather, "id = ?", DecodedID).Error; err != nil {
		session.Set("ERROR_WEATHER", "Weather not found")
		c.Redirect(http.StatusFound, "home/weather/")
		return
	}

	c.HTML(http.StatusOK, "detail_weather.html", gin.H{
		"title": "Detail Weather",
		"menu":  menu,
		"data":  weather,
		"user":  userSession,
	})
}

// GET /weather/edit/:hashid
func EditWeatherByID(c *gin.Context) {
	id := c.Param("id")
	DecodedID, _ := utils.DecodeID(id)
	session := sessions.Default(c)
	userSession := middlewares.GetSessionUser(c)

	session.Set("USERNAME_SESSION", userSession.Username)
	session.Save()

	menu, _ := utils.GetMenuSubmenu(c)
	var weather models.Weather

	if err := config.DB.First(&weather, "id = ?", DecodedID).Error; err != nil {
		session.Set("ERROR_WEATHER", "Weather not found")
		c.Redirect(http.StatusFound, "home/weather/")
		return
	}

	c.HTML(http.StatusOK, "edit_weather.html", gin.H{
		"title":     "Edit Weather",
		"menu":      menu,
		"data":      weather,
		"csrfToken": csrf.GetToken(c),
		"user":      userSession,
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

// POST /weather/:id
func UpdateWeather(c *gin.Context) {
	session := sessions.Default(c)
	id := c.Param("id")

	// Step 1: Define a temporary form struct
	type WeatherFormInput struct {
		City          string  `form:"City"`
		Summary       string  `form:"Summary"`
		Temperature   float64 `form:"Temperature"`
		WindSpeed     float64 `form:"WindSpeed"`
		WindAngle     float64 `form:"WindAngle"`
		WindDirection string  `form:"WindDirection"`
	}

	// Step 2: Bind form data
	var input WeatherFormInput
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", "Invalid input: "+err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	// Step 3: Find the weather record
	var weather models.Weather
	if err := config.DB.First(&weather, id).Error; err != nil {
		session.Set("ERROR", "Weather not found")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	// Step 4: Update fields manually
	weather.City = input.City
	weather.Summary = input.Summary
	weather.Temperature = input.Temperature
	weather.WindSpeed = input.WindSpeed
	weather.WindAngle = input.WindAngle
	weather.WindDirection = input.WindDirection

	// Step 5: Save to DB
	if err := config.DB.Save(&weather).Error; err != nil {
		session.Set("ERROR_UPDATE", "Failed to update weather: "+err.Error())
	} else {
		session.Set("SUCCESS_UPDATE", "Weather updated successfully")
	}
	session.Save()

	// Step 6: Redirect
	c.Redirect(http.StatusFound, "/home/weather")
}

// DELETE /weather/del:id
func DeleteWeather(c *gin.Context) {
	id := c.Param("id")
	DecodedID, _ := utils.DecodeID(id)
	session := sessions.Default(c)

	// Hard delete (bypass GORM's soft delete)
	if err := config.DB.Unscoped().Delete(&models.Weather{}, DecodedID).Error; err != nil {
		session.Set("DELETE_ERROR", "Delete failed")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	session.Set("DELETE_SUCCESS", "Weather successfully deleted")
	session.Save()
	c.Redirect(http.StatusFound, "/home/weather")
}
