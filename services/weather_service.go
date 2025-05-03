package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
			Direction string  `json:"dir"`
		} `json:"wind"`
	} `json:"current"`
}

func FetchAndUpdateWeather(location models.Location) {
	// Log START
	log.Printf("[CRON] - [START] - City: %s - Latitude: %.5f - Longitude: %.5f", location.City, location.Latitude, location.Longitude)

	apiKey := os.Getenv("METEO_API_KEY")
	if apiKey == "" {
		log.Println("[CRON] - [FAILED] - METEO_API_KEY tidak ditemukan di .env")
		return
	}

	url := fmt.Sprintf(
		"https://www.meteosource.com/api/v1/free/point?lat=%f&lon=%f&sections=current&timezone=auto&language=en&units=metric&key=%s",
		location.Latitude, location.Longitude, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("[CRON] - [FAILED] - City: %s - Error ambil data cuaca: %v", location.City, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data MeteoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("[CRON] - [FAILED] - City: %s - Error decode JSON: %v", location.City, err)
		return
	}

	// Cari apakah data sudah ada
	weather := models.Weather{
		LocationID:    location.ID,
		City:          location.City,
		Latitude:      location.Latitude,
		Longitude:     location.Longitude,
		Summary:       data.Current.Summary,
		Temperature:   data.Current.Temperature,
		WindSpeed:     data.Current.Wind.Speed,
		WindAngle:     data.Current.Wind.Angle,
		WindDirection: data.Current.Wind.Direction,
	}

	// ✅ Auto: jika ada → UPDATE, jika tidak → INSERT
	if err := config.DB.
		Where("location_id = ? AND city = ?", location.ID, location.City).
		Assign(weather).
		FirstOrCreate(&weather).Error; err != nil {
		log.Printf("[CRON] - [FAILED] - City: %s - Gagal simpan/update data: %v", location.City, err)
		return
	}

	// Delay simulasi async (opsional)
	time.Sleep(1 * time.Second)

	// ✅ Log Sukses
	log.Printf("[CRON] - [SUCCESS] - City: %s - Latitude: %.5f - Longitude: %.5f", location.City, location.Latitude, location.Longitude)
}
