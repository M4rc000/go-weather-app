package routes

import (
	"github.com/gin-gonic/gin"
	"weather-app/controllers"
	"weather-app/middlewares"
)

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.POST("/location", controllers.SaveLocation)
		
		v1.GET("/weather", controllers.GetAllWeather)
		v1.GET("/weather/:id", controllers.GetWeatherByID)
		v1.POST("/weather", controllers.CreateWeather)
		v1.PUT("/weather/:id", controllers.UpdateWeather)
		v1.DELETE("/weather/:id", controllers.DeleteWeather)

		v1.GET("/forecast/daily/:city", controllers.GetDailyForecast)
		v1.GET("/forecast/hourly/:city", controllers.GetHourlyForecast)

	}
}
