package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
	"net/http"
	"weather-app/config"
	"weather-app/models"
	"weather-app/utils"
)

func Register(c *gin.Context) {

	err := utils.FlashMessage(c, "ERROR")
	failedRegister := utils.FlashMessage(c, "FAILED_REGISTER")
	errorInputData := utils.FlashMessage(c, "ERROR_INPUTDATA")
	errorName := utils.FlashMessage(c, "ERROR_NAME")
	errorUsername := utils.FlashMessage(c, "ERROR_USERNAME")
	errorPassword := utils.FlashMessage(c, "ERROR_PASSWORD")
	duplicateUsername := utils.FlashMessage(c, "DUPLICATE_USERNAME")

	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":             "Register",
		"csrfToken":         csrf.GetToken(c),
		"failedRegister":    failedRegister,
		"errorName":         errorName,
		"errorUsername":     errorUsername,
		"error":             err,
		"errorPassword":     errorPassword,
		"errorInputData":    errorInputData,
		"duplicateUsername": duplicateUsername,
	})
}

func StoreRegister(c *gin.Context) {
	session := sessions.Default(c)
	var users models.User

	// 1. Bind input
	if err := c.ShouldBind(&users); err != nil {
		session.Set("ERROR", "Failed to bind input")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// 2. Validate input
	validate := validator.New()
	if err := validate.Struct(users); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range validationErrs {
				switch fe.Field() {
				case "Name":
					session.Set("ERROR_NAME", "Name is required")
				case "Username":
					session.Set("ERROR_USERNAME", "Username is required")
				case "Password":
					switch fe.Tag() {
					case "required":
						session.Set("ERROR_PASSWORD", "Password is required")
					case "min":
						session.Set("ERROR_PASSWORD", "Password must be at least 6 characters")
					}
				}
			}
		} else {
			session.Set("ERROR", "Invalid input")
		}
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// 3. Check duplicate
	var existing models.User
	if err := config.DB.Where("username = ?", users.Username).First(&existing).Error; err == nil {
		session.Set("DUPLICATE_USERNAME", "Username already exists")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	// 4. Hash password and save
	hashedPassword, _ := utils.HashPassword(users.Password)
	user := models.User{Username: users.Username, Password: hashedPassword}

	if err := config.DB.Create(&user).Error; err != nil {
		session.Set("FAILED_REGISTER", err.Error())
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	session.Set("SUCCESS_REGISTER", "Registration successfully")
	session.Save()
	c.Redirect(http.StatusFound, "/auth")
}

func Login(c *gin.Context) {
	successRegister := utils.FlashMessage(c, "SUCCESS_REGISTER")
	errorUsername := utils.FlashMessage(c, "ERROR_USERNAME")
	errorPassword := utils.FlashMessage(c, "ERROR_PASSWORD")
	loginError := utils.FlashMessage(c, "ERROR")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":           "Login",
		"csrfToken":       csrf.GetToken(c),
		"successRegister": successRegister,
		"loginError":      loginError,
		"errorUsername":   errorUsername,
		"errorPassword":   errorPassword,
	})
}

func StoreLogin(c *gin.Context) {
	session := sessions.Default(c)

	var input struct {
		gorm.Model
		Username string `gorm:"unique" form:"Username" validate:"required"`
		Password string `form:"Password" validate:"required"`
	}

	// Bind form input
	if err := c.ShouldBind(&input); err != nil {
		session.Set("ERROR", "Failed to bind input")
		session.Save()
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range validationErrs {
				switch fe.Field() {
				case "Username":
					session.Set("ERROR_USERNAME", "Username is required")
					session.Save()
				case "Password":
					switch fe.Tag() {
					case "required":
						session.Set("ERROR_PASSWORD", "Password is required")
						session.Save()
					}
				}
			}
		} else {
			session.Set("ERROR", "Invalid input")
			session.Save()
		}
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			session.Set("ERROR", "Invalid credentials")
			session.Save()
			c.Redirect(http.StatusFound, "/auth")
			return
		}
	}

	// Validate password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		session.Set("ERROR", "Invalid credentials")
		session.Save()
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	// Generate JWT Token
	tokenString, err := utils.GenerateToken(user.ID)
	if err != nil {
		session.Set("ERROR", "Login failed")
		session.Save()
		c.Redirect(http.StatusFound, "/auth")
		return
	}

	usernameSession := user.Username
	session.Set("USERID_SESSION", fmt.Sprintf("%d", user.ID))
	session.Set("USERNAME_SESSION", usernameSession)
	session.Set("JWT_TOKEN", tokenString)
	session.Set("SUCCESS_LOGIN", "Login successfully")
	session.Save()

	c.Redirect(http.StatusFound, "/home/get-location")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	session.Save()

	c.Redirect(http.StatusFound, "/auth")
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "NotFound.html", gin.H{
		"title": "Not Found",
	})
}
