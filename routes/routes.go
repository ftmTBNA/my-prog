package routes

import (
	"my-prog/handlers"
	"my-prog/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// مسیرهای عمومی
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.POST("/signup", handlers.SignupHandler)
	r.POST("/login", handlers.LoginHandler)


	auth := r.Group("/protected")
	auth.Use(middleware.AuthMiddleware())  
	{
		auth.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is a protected route!"})
		})
	}

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())  
    admin.Use(middleware.RequireRole("admin"))
    {
        admin.GET("/dashboard", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "Welcome to the admin dashboard"})
        })
    }

	return r
}
