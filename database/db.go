package database

import (
	"fmt"
	"os"

	"github.com/sing3demons/app/v2/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

//InitDB - connenct database
func InitDB() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable TimeZone=Asia/Bangkok", dbHost, user, password, dbName, db_port)
	// dsn := "host=localhost user=postgres password=12345678 dbname=product_shop port=30001 sslmode=disable TimeZone=Asia/Bangkok"
	database, err := gorm.Open(postgres.Open(dsn))
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
