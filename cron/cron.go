package cron

import (
	"context"
	"fmt"
	"log"
	"time"
	"weather-app/config"
	"weather-app/logbuffer"
	"weather-app/models"
	"weather-app/services"
)

var cancel context.CancelFunc
var running bool

func StartScheduler() {
	if running {
		log.Println("[CRON] Scheduler already running")
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancel = cancelFunc
	running = true

	log.Println("[CRON] Scheduler started")

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[CRON] Scheduler stopped")
				running = false
				return
			default:
				var locations []models.Location
				config.DB.Find(&locations)

				now := time.Now().Format("2006/01/02 15:04:05")
				log.Printf("[CRON] Start Schedule Job Cron at %s - Total Schedule Job: %d", now, len(locations))
				logbuffer.AddLog(fmt.Sprintf("[CRON] Start at %s", now))

				for _, loc := range locations {
					go services.FetchAndUpdateWeather(loc)
				}

				// Delay 1 Menit sebelum menjalankan request lagi
				time.Sleep(1 * time.Minute)
			}
		}
	}()
}

func StopScheduler() {
	if cancel != nil {
		cancel()
	}
}
