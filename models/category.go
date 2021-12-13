package models

import (
	"gorm.io/gorm"
)

//Category model
type Category struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Product []Product
}
