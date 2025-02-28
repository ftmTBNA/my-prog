package models

import "gorm.io/gorm"

type User struct{
	gorm.Model
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
	Role     string `gorm:"default:'user'"` // Default role is "user"
	// Password string `json:"password" gorm:"not null"`
}

func (user *User) CreateUser(db *gorm.DB) error {
	err := db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil

	// you can also use: return db.Create(user).Error
}

func FindByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}