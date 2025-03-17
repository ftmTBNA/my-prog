package middleware

import (
	"fmt"
	"my-prog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken, exists := c.Get("userToken")
		fmt.Println("userToken:", userToken)
		fmt.Println("exists:", exists)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// توی Go وقتی یه interface{} یا any داری (یا هر چیزی که تایپ مشخصی نداره)،
		// باید اونو به نوع واقعی‌ش تبدیل کنی تا بتونی باهاش کار کنی
		// مثلا اینجا خروجی تابع get از نوع any  هست و توی خط زیر به چیزی که میخواسته تبدیلش کرده
		tokenData, ok := userToken.(*utils.TokenData)
		if !ok || tokenData == nil {
			fmt.Println("Failed type assertion or tokenData is nil")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
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
