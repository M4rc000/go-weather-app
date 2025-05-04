package controllers

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
	"log"
	"net/http"
	"weather-app/config"
	"weather-app/models"
	"weather-app/utils"
)

func Register(c *gin.Context) {
	session := sessions.Default(c)

	// RETRIEVE FLASH MESSAGES
	err := session.Get("ERROR")
	failedRegister := session.Get("FAILED_REGISTER")
	errorInputData := session.Get("ERROR_INPUTDATA")
	errorName := session.Get("ERROR_NAME")
	errorUsername := session.Get("ERROR_USERNAME")
	errorPassword := session.Get("ERROR_PASSWORD")
	duplicateUsername := session.Get("DUPLICATE_USERNAME")

	//Clear flash messages (so the message disappears after refresh)
	session.Delete("FAILED_REGISTER")
	session.Delete("ERROR_NAME")
	session.Delete("ERROR_USERNAME")
	session.Delete("ERROR_PASSWORD")
	session.Delete("DUPLICATE_USERNAME")

	session.Save()

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
	session := sessions.Default(c)
	err := session.Get("ERROR")
	errorUsername := session.Get("ERROR_USERNAME")
	errorPassword := session.Get("ERROR_PASSWORD")
	successRegister := session.Get("SUCCESS_REGISTER")
	loginError := session.Get("LOGIN_ERROR")

	session.Delete("LOGIN_ERROR")
	session.Delete("SUCCESS_REGISTER")

	session.Save()

	c.HTML(http.StatusFound, "login.html", gin.H{
		"title":           "Login",
		"csrfToken":       csrf.GetToken(c),
		"error":           err,
		"errorUsername":   errorUsername,
		"errorPassword":   errorPassword,
		"successRegister": successRegister,
		"loginError":      loginError,
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
				case "Password":
					switch fe.Tag() {
					case "required":
						session.Set("ERROR_PASSWORD", "Password is required")
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

	// Set token in secure HttpOnly cookie
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true) // secure=true if HTTPS

	// Optional: Store a success flash
	session.Delete("USERID")
	session.Delete("USERNAME")

	session.Set("USERID", user.ID)
	session.Set("USERNAME", user.Username)
	session.Set("LOGIN_SUCCESS", "Login successfully")
	if err := session.Save(); err != nil {
		log.Println("Failed to save session:", err)
	}

	c.Redirect(http.StatusFound, "/home/get-location")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)

	// Hapus semua data dari session
	session.Clear()
	session.Save()

	// Hapus token dari cookie (jika digunakan)
	c.SetCookie("token", "", -1, "/", "", false, true) // expire cookie

	// Redirect ke halaman login atau halaman utama
	c.Redirect(http.StatusFound, "/auth")
}
