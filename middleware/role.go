package middleware

import (
	"my-prog/utils"
    "net/http"
	"fmt"

	"github.com/gin-gonic/gin"

)


func RequireRole(role string) gin.HandlerFunc  {
	return func(c *gin.Context) {
		userToken, exists := c.Get("userToken")
		fmt.Println("userToken:", userToken) 
		fmt.Println("exists:", exists) 
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenData, ok := userToken.(*utils.TokenData)
		fmt.Println("tokenData:", tokenData) 
		fmt.Println("tokenData.Role:", tokenData.Role) 
		fmt.Println("role:", role) 
		fmt.Println("ok:", ok) 
        if !ok || tokenData.Role != role {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}