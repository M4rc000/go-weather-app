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
			var locations []models.Location
			config.DB.Find(&locations)

			now := time.Now().Format("2006/01/02 15:04:05")
			log.Printf("[CRON] Start Schedule Job Cron at %s - Total Schedule Job: %d", now, len(locations))

			for _, loc := range locations {
				go services.FetchAndUpdateWeather(loc)
			}
			// Delay 1 Menit sebelum menjalankan request lagi
			time.Sleep(1 * time.Minute)
		}
	}()
}
