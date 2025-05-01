package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"weather-app/config"
	"weather-app/cron"
	"weather-app/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	cron.StartScheduler()
	app := gin.Default()
	routes.SetupRoutes(app)

	log.Println("App running on Port 3000")
	app.Run(":3000")
}
