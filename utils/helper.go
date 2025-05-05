package utils

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/speps/go-hashids"
	"os"
	"strings"
	"time"
	"unicode"
)

func Proper(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func GetMenuSubmenu(c *gin.Context) (menu, submenu string) {
	URL := strings.Split(c.Request.URL.Path, "/")
	menu = Proper(URL[1])
	submenu = Proper(URL[2])
	return menu, submenu
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// pastikan menggunakan algoritma yang sesuai
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

// GENERATE TOKEN SAAT LOGIN
func GenerateToken(userID uint) (string, error) {
	// Claims berisi informasi yang ingin kamu simpan di token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // token berlaku 1 jam
		"iat":     time.Now().Unix(),                    // issued at
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret key
	secret := []byte(os.Getenv("JWT_SECRET")) // misalnya "mysecretkey"
	return token.SignedString(secret)
}

func FlashMessage(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	val := session.Get(key)
	session.Delete(key)
	session.Save()
	return val
}

func EncodeID(id int) string {
	hd := hashids.NewData()
	hd.Salt = "your-secure-salt" // Use a strong, secret salt
	hd.MinLength = 6             // Optional: Min length of encoded string
	h, _ := hashids.NewWithData(hd)

	e, _ := h.Encode([]int{id})
	return e
}

func DecodeID(encoded string) (int, error) {
	hd := hashids.NewData()
	hd.Salt = "your-secure-salt" // Same salt as above!
	hd.MinLength = 6
	h, _ := hashids.NewWithData(hd)

	ids, err := h.DecodeWithError(encoded)
	if err != nil || len(ids) == 0 {
		return 0, fmt.Errorf("invalid ID")
	}

	return ids[0], nil
}
