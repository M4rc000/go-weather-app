package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"weather-app/config"
	"weather-app/models"
)

type MeteoResponse struct {
	Current struct {
		Summary     string  `json:"summary"`
		Temperature float64 `json:"temperature"`
		Wind        struct {
			Speed     float64 `json:"speed"`
			Angle     float64 `json:"angle"`
			Direction string  `json:"direction"`
		} `json:"wind"`
	} `json:"current"`
}

func FetchAndUpdateWeather(location models.Location) {
	apiKey := os.Getenv("METEO_API_KEY")
	if apiKey == "" {
		log.Println("⚠️  METEO_API_KEY not set")
		return
	}

	url := fmt.Sprintf(
		"https://www.meteosource.com/api/v1/free/point?lat=%f&lon=%f&sections=current&timezone=auto&language=en&units=metric&key=%s",
		location.Latitude,
		location.Longitude,
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("🌩️ Failed to fetch weather for %s: %v\n", location.City, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("❌ API error for %s: %s\n", location.City, string(body))
		return
	}

	var data MeteoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("❌ Failed to decode JSON:", err)
		return
	}

	// Update ke tabel Weather
	weather := models.Weather{
		LocationID:    location.ID,
		Summary:       data.Current.Summary,
		Temp:          data.Current.Temperature,
		WindSpeed:     data.Current.Wind.Speed,
		WindAngle:     data.Current.Wind.Angle,
		WindDirection: data.Current.Wind.Direction,
	}

	if err := config.DB.Create(&weather).Error; err != nil {
		log.Println("❌ Failed to save weather:", err)
		return
	}

	log.Printf("✅ Weather updated for %s\n", location.City)
}
