package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"weather-app/utils"

	"github.com/gin-gonic/gin"
)

func GetDailyForecast(c *gin.Context) {
	city := c.Param("city")
	cacheKey := "daily:" + city

	// ✅ Check from cache
	if cached, found := utils.GetCache(cacheKey); found {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success (cached)",
			"message": "Daily forecast (cached)",
			"data":    cached,
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

	// ✅ Set to cache
	utils.SetCache(cacheKey, result["daily"], 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daily forecast for " + city,
		"data":    result["daily"],
	})
}

func GetHourlyForecast(c *gin.Context) {
	city := c.Param("city")
	cacheKey := "hourly:" + city

	// ✅ Check from cache
	if cached, found := utils.GetCache(cacheKey); found {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success (cached)",
			"message": "Hourly forecast (cached)",
			"data":    cached,
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

	// ✅ Set to cache
	utils.SetCache(cacheKey, result["hourly"], 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Hourly forecast for " + city,
		"data":    result["hourly"],
	})
}
