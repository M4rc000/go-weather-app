package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"log"
	"net/http"
	"weather-app/config"
	"weather-app/models"
	"weather-app/utils"
)

func GetLocation(c *gin.Context) {
	session := sessions.Default(c)
	menu, submenu := utils.GetMenuSubmenu(c)

	username := session.Get("USERNAME")

	c.HTML(http.StatusFound, "get_location_data.html", gin.H{
		"title":     "Get Location",
		"csrfToken": csrf.GetToken(c),
		"menu":      menu,
		"submenu":   submenu,
		"user":      username,
	})
}

func SearchLocation(c *gin.Context) {
	var input struct {
		City string `form:"city" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "City is required",
		})
		return
	}

	var locations []models.Location

	// Search city with LIKE
	if err := config.DB.Where("city LIKE ?", "%"+input.City+"%").Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Terjadi kesalahan saat pencarian",
		})
		return
	}

	// No result
	if len(locations) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "empty",
			"message": "Data tidak ditemukan",
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"locations": locations,
	})
}

func SaveLocation(c *gin.Context) {
	session := sessions.Default(c)

	var input struct {
		City      string  `form:"City" validate:"required"`
		Latitude  float64 `form:"Latitude" validate:"required"`
		Longitude float64 `form:"Longitude" validate:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", err.Error())
		session.Save()
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range validationErrs {
				switch fe.Field() {
				case "City":
					session.Set("ERROR_CITY", "City is required")
				case "Latitude":
					session.Set("ERROR_LATITUDE", "Latitude is required")
				case "Longitude":
					session.Set("ERROR_LONGITUDE", "Longitude is required")
				}
			}
		} else {
			session.Set("ERROR", "Invalid input")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// Cek duplicate berdasarkan city + latitude + longitude
	var existing models.Location
	if err := config.DB.Where("city = ? AND latitude = ? AND longitude = ?",
		input.City, input.Latitude, input.Longitude).First(&existing).Error; err == nil {
		log.Printf("[DUPLIKAT] Lokasi %s sudah ada di database. ID: %d", existing.City, existing.ID)
		session.Set("DUPLICATE_LOCATION", "Location already exists")
		session.Save()
		return
	}

	// Simpan data
	location := models.Location{
		City:      input.City,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Weather: models.Weather{
			Summary:       "-",
			Temperature:   0,
			WindSpeed:     0,
			WindAngle:     0,
			WindDirection: "-",
		},
	}

	if err := config.DB.Create(&location).Error; err != nil {
		session.Set("FAILED_SAVE", "Failed to save location")
		session.Save()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"locations": []gin.H{
			{
				"ID":              location.ID,
				"CreatedAt":       location.CreatedAt,
				"UpdatedAt":       location.UpdatedAt,
				"DeletedAt":       location.DeletedAt,
				"city":            location.City,
				"latitude":        location.Latitude,
				"longitude":       location.Longitude,
				"weather_summary": location.Weather.Summary,
				"temperature":     location.Weather.Temperature,
				"wind_speed":      location.Weather.WindSpeed,
				"wind_angle":      location.Weather.WindAngle,
				"wind_direction":  location.Weather.WindDirection,
			},
		},
	})

}
