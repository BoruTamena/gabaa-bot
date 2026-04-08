package middleware

import (
	"net/http"
	"strings"



	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AuthMiddleware struct {
	tele platform.Telegram
}

func NewAuthMiddleware(tele platform.Telegram) *AuthMiddleware {
	return &AuthMiddleware{tele: tele}
}

func (m *AuthMiddleware) TelegramAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := c.GetHeader("X-Telegram-Init-Data")
		if initData == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing telegram init data"})
			c.Abort()
			return
		}

		valid, err := m.tele.ValidateInitData(initData)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid telegram init data"})
			c.Abort()
			return
		}

		// In a real app, you'd parse user_id from initData here.
		// For now, we assume it's passed or extracted later.
		c.Next()
	}
}

func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int64(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}
