package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func Guest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		tokenString, ok := session.Get("JWT_TOKEN").(string)
		if !ok || tokenString == "" {
			// Tidak ada token, user dianggap guest → boleh lanjut ke login/register
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err == nil && token.Valid {
			// Sudah login → redirect ke home
			c.Redirect(http.StatusFound, "/home")
			return
		}

		// Token tidak valid → biarkan akses auth page
		c.Next()
	}
}
