package middlewares

import (
	"net/http"
	"weather-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token") // atau ambil dari header jika pakai Authorization

		if err != nil || tokenString == "" {
			c.Redirect(http.StatusFound, "/auth")
			return
		}

		token, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/auth")
			return
		}

		// Ambil claims (misalnya userID atau username) dan masukkan ke context
		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"])
		c.Next()
	}
}
