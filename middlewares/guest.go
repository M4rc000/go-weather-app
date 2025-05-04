package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func Guest() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			// No token, proceed to login/register
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err == nil && token.Valid {
			// Already logged in, redirect away from auth pages
			c.Redirect(http.StatusFound, "/home")
			c.Abort()
			return
		}

		// Invalid token, allow access to auth pages
		c.Next()
	}
}
