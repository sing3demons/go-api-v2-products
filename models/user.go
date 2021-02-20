package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User model
type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Name     string `gorm:"not null"`
	Avatar   string
	Role     string `gorm:"default:'Member';not null"`
}

//GenerateEncryptedPassword - hash password
func (u *User) GenerateEncryptedPassword() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return string(hash)
}
