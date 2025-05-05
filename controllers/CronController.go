package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"weather-app/cron"
	"weather-app/logbuffer"
	"weather-app/middlewares"
	"weather-app/utils"
)

func CronJob(c *gin.Context) {
	session := sessions.Default(c)
	userSession := middlewares.GetSessionUser(c)
	session.Set("USERNAME_SESSION", userSession.Username)
	session.Save()

	menu, _ := utils.GetMenuSubmenu(c)
	c.HTML(http.StatusOK, "cron_logs.html", gin.H{
		"title": "Cron Job Scheduler",
		"menu":  menu,
		"user":  userSession,
	})
}

func GetCronLogs(c *gin.Context) {
	logs := logbuffer.GetLogs()
	c.JSON(200, gin.H{"logs": logs})
}

func RunCron(c *gin.Context) {
	go cron.StartScheduler() // run it in a goroutine so it doesn't block
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cron scheduler started.",
	})
}

func StopCron(c *gin.Context) {
	cron.StopScheduler()
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cron scheduler stopped.",
	})
}
