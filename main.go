package main

import (
	"log"
	"my-prog/database"
	"my-prog/routes"
	"my-prog/utils"

	"github.com/joho/godotenv"
)

// func init() {
//     // Load environment variables from .env file
//     if err := godotenv.Load(); err != nil {
//         log.Println("No .env file found")
//     }
// }


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Println(".env file loaded successfully")
	}
	utils.Init()
}

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