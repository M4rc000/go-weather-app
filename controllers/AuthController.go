package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
	"weather-app/config"
	"weather-app/models"
	"weather-app/utils"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	//  CEK DUPLICATE USERNAME
	var existing models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Username already exists"})
		return
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	user := models.User{Username: input.Username, Password: hashedPassword}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User has successfully registered"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Customisasi Pesan Error
		var validation_error validator.ValidationErrors
		if errors.As(err, &validation_error) {
			field := validation_error[0].Field()
			tag := validation_error[0].Tag()

			var msg string
			switch field {
			case "Username":
				msg = "Username is required"
			case "Password":
				if tag == "required" {
					msg = "Password is required"
				} else if tag == "min" {
					msg = "Password must be at least 6 characters"
				}
			default:
				msg = "Invalid input"
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": msg,
			})
			return
		}
		// Fallback error
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request",
		})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successfully", "token": tokenString})
}
