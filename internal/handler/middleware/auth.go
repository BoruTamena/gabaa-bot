package middleware

import (
	"net/http"
	"strings"

	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
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
			appErr := errorx.New(errorx.ErrUnauthorized, "Missing telegram init data", http.StatusUnauthorized)
			c.JSON(appErr.Status, appErr)
			c.Abort()
			return
		}

		valid, err := m.tele.ValidateInitData(initData)
		if err != nil || !valid {
			appErr := errorx.New(errorx.ErrUnauthorized, "Invalid telegram init data", http.StatusUnauthorized)
			c.JSON(appErr.Status, appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := errorx.New(errorx.ErrUnauthorized, "Missing authorization header", http.StatusUnauthorized)
			c.JSON(appErr.Status, appErr)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil || !token.Valid {
			appErr := errorx.New(errorx.ErrUnauthorized, "Invalid or expired token", http.StatusUnauthorized)
			c.JSON(appErr.Status, appErr)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			appErr := errorx.New(errorx.ErrUnauthorized, "Invalid token claims", http.StatusUnauthorized)
			c.JSON(appErr.Status, appErr)
			c.Abort()
			return
		}

		c.Set("user_id", int64(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))
		if storeID, ok := claims["store_id"].(float64); ok {
			c.Set("store_id", int64(storeID))
		}
		
		c.Next()
	}
}
