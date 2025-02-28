package utils

import (
	"errors"
	"regexp"
	"time"
	"os"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
// "os"

var secretKey string

func init() {
	secretKey = os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET is not set")
	} else {
		log.Println("JWT_SECRET loaded successfully")
	}
}


type TokenData struct {
    UserID uint   `json:"userId"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// **هش کردن رمز عبور**
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// **بررسی صحت رمز عبور**
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// **تولید توکن JWT**
// func GenerateToken(userID uint) (string, error) {
// 	secretKey := []byte(os.Getenv("JWT_SECRET" )) // مقدار متغیر محیطی
// 	fmt.Println("JWT_SECRET:", secretKey) 
// 	if len(secretKey) == 0 {
// 		return "", errors.New("JWT_SECRET is not set")
// 	}

// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 ساعت اعتبار
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(secretKey)
// }

// const secretKey = "supersecret"
	

func GenerateToken(email string, userID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	if secretKey == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	log.Println("JWT_SECRET loaded successfully")


	return token.SignedString([]byte(secretKey))
}

// **بررسی اعتبار توکن JWT**
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// secretKey := []byte(os.Getenv("JWT_SECRET"))
	secretKey := []byte(secretKey)
	if len(secretKey) == 0 {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	return token, err
}

// **بررسی صحت ایمیل**
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
