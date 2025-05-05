package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	csrf "github.com/utrack/gin-csrf"
	"io"
	"net/http"
	"os"
	"time"
	"weather-app/middlewares"
	"weather-app/utils"

	"github.com/gin-gonic/gin"
)

func GetForecast(c *gin.Context) {
	session := sessions.Default(c)
	userSession := middlewares.GetSessionUser(c)
	session.Set("USERNAME_SESSION", userSession.Username)
	session.Save()
	menu, _ := utils.GetMenuSubmenu(c)

	c.HTML(http.StatusFound, "show_forecast.html", gin.H{
		"title":     "Forecast",
		"menu":      menu,
		"csrfToken": csrf.GetToken(c),
		"user":      userSession,
	})
}

func GetDailyForecast(c *gin.Context) {
	city := c.Param("city")
	cacheKey := "daily:" + city

	// ✅ Check from cache and ensure it's an array
	if cached, found := utils.GetCache(cacheKey); found {
		var dailyData []interface{}
		if data, ok := cached.([]interface{}); ok {
			dailyData = data
		} else if cached != nil {
			dailyData = []interface{}{cached}
		} else {
			dailyData = []interface{}{}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "success (cached)",
			"message": "Daily forecast (cached)",
			"data":    dailyData,
		})
		return
	}

	apiKey := os.Getenv("METEO_API_KEY")
	url := fmt.Sprintf(
		"https://www.meteosource.com/api/v1/free/point?place_id=%s&sections=daily&timezone=auto&language=en&units=metric&key=%s",
		city,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to fetch daily forecast"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// ✅ Convert to a slice if necessary
	var dailyData []interface{}
	if data, ok := result["daily"].([]interface{}); ok {
		dailyData = data
	} else if result["daily"] != nil {
		dailyData = []interface{}{result["daily"]}
	} else {
		dailyData = []interface{}{}
	}

	// ✅ Set to cache
	utils.SetCache(cacheKey, dailyData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daily forecast for " + city,
		"data":    dailyData,
	})
}

func GetHourlyForecast(c *gin.Context) {
	city := c.Param("city")
	cacheKey := "hourly:" + city

	// ✅ Check from cache and ensure it's an array
	if cached, found := utils.GetCache(cacheKey); found {
		var hourlyData []interface{}
		if data, ok := cached.([]interface{}); ok {
			hourlyData = data
		} else if cached != nil {
			hourlyData = []interface{}{cached}
		} else {
			hourlyData = []interface{}{}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "success (cached)",
			"message": "Hourly forecast (cached)",
			"data":    hourlyData,
		})
		return
	}

	apiKey := os.Getenv("METEO_API_KEY")
	url := fmt.Sprintf(
		"https://www.meteosource.com/api/v1/free/point?place_id=%s&sections=hourly&timezone=auto&language=en&units=metric&key=%s",
		city,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to fetch hourly forecast"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// ✅ Convert to a slice if necessary
	var hourlyData []interface{}
	if data, ok := result["hourly"].([]interface{}); ok {
		hourlyData = data
	} else if result["hourly"] != nil {
		hourlyData = []interface{}{result["hourly"]}
	} else {
		hourlyData = []interface{}{}
	}

	// ✅ Set to cache
	utils.SetCache(cacheKey, hourlyData, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Hourly forecast for " + city,
		"data":    hourlyData,
	})
}
