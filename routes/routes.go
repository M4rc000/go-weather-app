package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"weather-app/controllers"
	"weather-app/middlewares"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/auth")
	})

	// Handle not found route
	r.NoRoute(controllers.NotFound)

	auth := r.Group("/auth")
	auth.Use(middlewares.Guest())
	{
		auth.GET("/register", controllers.Register)
		auth.POST("/register", controllers.StoreRegister)
		auth.GET("", controllers.Login)
		auth.POST("/login", controllers.StoreLogin)
	}

	home := r.Group("/home")
	home.Use(middlewares.Authenticate())
	{
		home.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/home/get-location")
		})
		home.GET("/get-location", controllers.GetLocation)
		home.POST("/get-location", controllers.SearchLocation)

		home.GET("/weather", controllers.GetAllWeather)
		home.GET("/weather/detail/:id", controllers.GetWeatherByID)
		home.GET("/weather/edit/:id", controllers.EditWeatherByID)
		home.POST("/weather", controllers.CreateWeather)
		home.POST("/weather/update/:id", controllers.UpdateWeather)
		home.GET("/weather/delete/:id", controllers.DeleteWeather)

		home.GET("/cron", controllers.CronJob)
		home.GET("/cron/logs", controllers.GetCronLogs)
		home.GET("/run-cron", controllers.RunCron)
		home.GET("/stop-cron", controllers.StopCron)

		home.GET("/forecast", controllers.GetForecast)

		home.GET("/logout", controllers.Logout)
	}

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.Authenticate())
	{
		v1.POST("/location", controllers.SaveLocation)
		v1.POST("/search-location", controllers.SearchLocation)

		v1.GET("/forecast/daily/:city", controllers.GetDailyForecast)
		v1.GET("/forecast/hourly/:city", controllers.GetHourlyForecast)

	}
}
