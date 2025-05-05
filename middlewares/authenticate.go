package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"weather-app/config"
	"weather-app/models"
	"weather-app/utils"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		tokenString, ok := session.Get("JWT_TOKEN").(string)
		if !ok || tokenString == "" {
			c.Redirect(http.StatusFound, "/auth")
			return
		}

		token, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/auth")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"])
		c.Next()
	}
}

func GetSessionUser(c *gin.Context) *models.User {
	Username := utils.FlashMessage(c, "USERNAME_SESSION")

	var user models.User
	if err := config.DB.Where("username = ?", Username).First(&user).Error; err != nil {
		c.HTML(http.StatusNotFound, "NotFound.html", gin.H{"error": "User not found"})
		return nil
	}

	return &user
}
