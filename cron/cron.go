package cron

import (
	"log"
	"time"
	"weather-app/config"
	"weather-app/models"
	"weather-app/services"
)

func StartScheduler() {
	go func() {
		for {
			log.Println("⏰ Running scheduled weather update...")

			var locations []models.Location
			if err := config.DB.Find(&locations).Error; err != nil {
				log.Println("❌ Error fetching locations:", err)
				time.Sleep(5 * time.Minute)
				continue
			}

			// Async: goroutine per lokasi
			for _, loc := range locations {
				go services.FetchAndUpdateWeather(loc)
			}

			// Batasi kecepatan sesuai MeteoSource plan: 10/min
			time.Sleep(1 * time.Minute)
		}
	}()
}
