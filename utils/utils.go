package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// "os"

var secretKey string
// secretKey = []byte(os.Getenv("JWT_SECRET"))

func Init() {

	secretKey = os.Getenv("JWT_SECRET")
	fmt.Println("secretKey:", secretKey)


	// secretKey = os.Getenv("JWT_SECRET")
	// if secretKey == "" {
	// 	log.Fatal("JWT_SECRET is not settttt")
	// } else {
	// 	log.Println("JWT_SECRET loaded successfully")
	// }
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

func GenerateToken(email string, userID uint, role string) (string, error) {
	claims := TokenData{
		Email:  email,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if secretKey == "" {
		log.Fatal("JWT_SECRET is not settt")
	}

	return token.SignedString([]byte(secretKey))
}

// **بررسی اعتبار توکن JWT**
func VerifyToken(tokenString string) (*jwt.Token, error) {
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET is not settt")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenData{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// **بررسی صحت ایمیل**
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
