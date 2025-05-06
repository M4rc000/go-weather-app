package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/go-playground/validator/v10"
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
	successCreate := utils.FlashMessage(c, "SUCCESS_CREATE")
	errorUpdate := utils.FlashMessage(c, "ERROR_UPDATE")
	deleteError := utils.FlashMessage(c, "DELETE_ERROR")
	duplicatedWeather := utils.FlashMessage(c, "DUPLICATE_WEATHER")

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

	c.HTML(http.StatusOK, "show_weather.html", gin.H{
		"title":             "Weather",
		"menu":              menu,
		"data":              DataWeathers,
		"csrfToken":         csrf.GetToken(c),
		"user":              userSession,
		"errorWeather":      errorWeather,
		"err":               err,
		"successUpdate":     successUpdate,
		"errorUpdate":       errorUpdate,
		"successCreate":     successCreate,
		"successDelete":     successDelete,
		"deleteError":       deleteError,
		"duplicatedWeather": duplicatedWeather,
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

// GET /weather/edit/:id
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

// POST /weather/create
func CreateWeather(c *gin.Context) {
	session := sessions.Default(c)

	var input struct {
		Latitude      float64 `form:"Latitude" gorm:"column:latitude"`
		Longitude     float64 `form:"Longitude" gorm:"column:longitude"`
		City          string  `form:"City" gorm:"column:city" gorm:"required"`
		Summary       string  `form:"Summary" gorm:"column:summary" gorm:"required"`
		Temperature   float64 `form:"Temperature" gorm:"column:temperature" gorm:"required"`
		WindSpeed     float64 `form:"WindSpeed" gorm:"column:wind_speed" gorm:"required"`
		WindAngle     float64 `form:"WindAngle" gorm:"column:wind_angle" gorm:"required"`
		WindDirection string  `form:"WindDirection" gorm:"column:wind_direction" gorm:"required"`
		LocationID    int     `gorm:"column:location_id"`
	}

	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", "All the field are required")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	// Validate Input
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range validationErrs {
				switch fe.Field() {
				case "City":
					session.Set("ERROR", "City is required")
				case "Summary":
					session.Set("ERROR", "Summary is required")
				case "Temperature":
					session.Set("ERROR", "Temperature is required")
				case "WindSpeed":
					session.Set("ERROR", "Wind Speed is required")
				case "WindAngle":
					session.Set("ERROR", "Wind Angle is required")
				case "WindDirection":
					session.Set("ERROR", "Wind Direction is required")
				}
			}
		} else {
			session.Set("ERROR", "Invalid input")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Cek duplicate weather data berdasarkan nama kota
	var existing models.Weather
	if err := config.DB.Where("city = ?", input.City).First(&existing).Error; err == nil {
		session.Set("DUPLICATE_WEATHER", "Data weather already exists")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	// Cari data lokasi dari table locations
	var loc models.Location
	if err := config.DB.Where("city ILIKE ?", "%"+input.City+"%").First(&loc).Error; err == nil {
		// Ubah input dengan data lokasi yang ditemukan
		input.LocationID = int(loc.ID)
		input.City = loc.City
		input.Latitude = loc.Latitude
		input.Longitude = loc.Longitude
	} else {
		input.LocationID = 0
	}

	// Simpan data weather ke tabel "weathers"
	if err := config.DB.Table("weathers").Create(&input).Error; err != nil {
		session.Set("ERROR", "Failed to save weather")
		session.Save()
		c.Redirect(http.StatusFound, "/home/weather")
		return
	}

	session.Set("SUCCESS_CREATE", "New weather successfully created")
	session.Save()
	c.Redirect(http.StatusFound, "/home/weather")
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
