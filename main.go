package main

import (
	"log"
	"my-prog/database"
	"my-prog/routes"
)
func main(){
	// r:=gin.Default()

	// r.GET("/users",handlers.GetUsers)
	// r.GET("/users",handlers.GetUsers)
	// r.GET("/users",handlers.GetUsers)

	database.ConnectDatabase()

	// تنظیم و اجرای روتر
	r := routes.SetupRouter()
	log.Println("🚀 Server is running on port 8080")

	r.Run(":8080")
}