package config

import (
	"app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

//InitDB - connenct database
func InitDB() {
	database, err := gorm.Open(sqlite.Open("./gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.User{})
	// database.Migrator().DropTable(&models.User{})

	db = database
}

//GetDB - return db
func GetDB() *gorm.DB {
	return db
}
