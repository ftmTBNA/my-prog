package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"my-prog/database"
	"my-prog/models"
	"my-prog/redis-caching"
	"my-prog/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// var users = []models.User{
// 	{ID: 1, Name: "Ali", Email: "ali@example.come"},
// 	{ID: 2, Name: "Sara", Email: "sara@example.come"},
// }

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	cacheKey := fmt.Sprintf("user:%d", id)
	cachedData, err := rediss.RedisClient.Get(rediss.Ctx, cacheKey).Result()
	if err == nil {
		log.Println("Cache hit...")
		if err := json.Unmarshal([]byte(cachedData), &user); err != nil {
			log.Println("Failed to unmarshal user: %v", err)
			return 
		}
		c.JSON(http.StatusOK, user)
	}

	if err != redis.Nil {
        log.Printf("Redis error: %v", err)
    }

	// var users []models.User


	err = database.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in DB"})
		return
	}

	log.Println("Cache miss! fetched from Postgres")

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("JSON Marshal error: %v", err)
	}else {
		err := rediss.RedisClient.Set(rediss.Ctx, cacheKey,userJson, 10 * time.Minute).Err()
		if err != nil {
            log.Printf("Failed to cache user: %v", err)
        } else {
            log.Println("User cached in Redis!")
        }
	}

	c.JSON(http.StatusOK, user)
}

// func CreateUser(c *gin.Context) {
// 	var NewUser models.User

// 	err := c.ShouldBindJSON(&NewUser)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Save data in db -------------------------------
// 	// METHOD 1

// 	// در این روش، ما از یک تابع جداگانه (CreateUser) که در پکیج models تعریف شده، استفاده می‌کنیم.
// 	// 🔹 این روش باعث می‌شه کد تمیزتر، قابل استفاده مجدد و ماژولار بشه.

// 	// METHOD 2
// 	// database.DB.Create(&NewUser)
// 	// در این روش، ما مستقیماً از gorm برای ذخیره کاربر جدید در دیتابیس استفاده می‌کنیم.
// 	// 🔹 مشکلات این روش:

// 	// اگر خطایی در ذخیره‌سازی رخ بده، نمی‌تونیم اون رو مدیریت کنیم مگر اینکه if err := database.DB.Create(...).Error رو دستی چک کنیم.
// 	// کد ما وابسته به دیتابیس (database.DB) در سطح handlers می‌شه، که باعث کاهش ماژولار بودن کد می‌شه.

// 	c.JSON(http.StatusCreated, NewUser)

// }

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// err := database.DB.First(&user, id).Error
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 	return
	// }

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(models.User{Name: user.Name, Email: user.Email})
    if result.Error != nil {
        log.Printf("Failed to update user in Postgres: %v", result.Error)
        return
    }

	// database.DB.Save(&user)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	err := database.DB.First(&user, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}




// SignupHandler مدیریت ثبت‌نام کاربران جدید
func SignupHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// بررسی معتبر بودن ایمیل
	if !utils.IsValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// هش کردن رمز عبور
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword


	// Default role is "user" unless set by an admin
    if user.Role == "" {
        user.Role = "user"
    }

	// ذخیره در دیتابیس
	if err := user.CreateUser(database.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// LoginHandler مدیریت ورود کاربران
func LoginHandler(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// جستجوی کاربر بر اساس ایمیل
	user, err := models.FindByEmail(database.DB, input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// بررسی رمز عبور
	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// تولید توکن
	fmt.Println("User ID:", user.ID) 
	                                  //input.Email
	token, err := utils.GenerateToken(user.Email, uint(user.ID), user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
