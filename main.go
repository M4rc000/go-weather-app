package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	csrf "github.com/utrack/gin-csrf"
	"log"
	"net/http"
	"os"
	"weather-app/config"
	"weather-app/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	//cron.StartScheduler()

	app := gin.Default()
	app.Static("/assets", "./assets")
	app.LoadHTMLGlob("views/**/*")

	// Set Cookie and Session Config
	store := cookie.NewStore([]byte("secret"))

	store.Options(sessions.Options{
		MaxAge:   3600,                 // Session lifetime in seconds (1 hour)
		Path:     "/",                  // Cookie available on all paths
		HttpOnly: true,                 // Prevent access from JavaScript (recommended)
		Secure:   false,                // Set to true in production if using HTTPS
		SameSite: http.SameSiteLaxMode, // Controls cross-site cookie behavior
	})

	app.Use(sessions.Sessions("weatherapp", store))

	// Set CRSF TOKEN
	app.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CRSF_TOKEN"),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	routes.SetupRoutes(app)

	log.Println("App running on Port 3000")
	errs := app.Run(":3000")
	if errs != nil {
		log.Println("Failed to running on port 3000")
	}
}
