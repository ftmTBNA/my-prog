package database

import (
	"log"
	"my-prog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {


// 	requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_NAME", "DB_HOST", "DB_PORT"}
// for _, v := range requiredVars {
//     if os.Getenv(v) == "" {
//         log.Fatalf("Environment variable %s is not set", v)
//     }
// }

	// dsn := "host=localhost user=postgres password=9912134 dbname=userdb port=5432 sslmode=disable"
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})   
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
