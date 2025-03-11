package middleware

import (
	"fmt"
	"my-prog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// "strings"
// AuthMiddleware بررسی توکن JWT برای احراز هویت
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}
		fmt.Println("token1: ", authHeader)
		// جدا کردن Bearer Token
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.VerifyToken(authHeader)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		
		// اینجا Claims رو بکش بیرون
		claims, ok := token.Claims.(*utils.TokenData)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		
		c.Set("userToken", claims)
		
		c.Next()
		
	}
}
