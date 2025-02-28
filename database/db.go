package database

import (
	"log"
	"my-prog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=9912134 dbname=userdb port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // ✅ استفاده از سینتکس جدید
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = database.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("AutoMigrate failed:", err)
    }

    DB = database
	
	log.Println("Database connected successfully!")


	
    // ایجاد جدول کاربران اگر وجود نداشته باشد
	// database.AutoMigrate(&models.User{})


	// برای اطمینان از صحت اتصال، یک instance از *sql.DB بگیر و تست کن
	// sqlDB, err := DB.DB()
	// if err != nil {
	// 	log.Fatal("Failed to get database instance:", err)
	// }

	// // تنظیمات کانکشن (اختیاری ولی توصیه‌شده)
	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// // sqlDB.SetConnMaxLifetime(0)

	// log.Println("Database connection is configured.")
}
