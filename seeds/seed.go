package seeds

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/app/v2/database"
	"github.com/sing3demons/app/v2/models"
)

func Load() {
	db := database.GetDB()

	var users []models.User
	db.Find(&users)
	if len(users) == 0 {
		fmt.Println("Creating users...")
		db.Migrator().DropTable(&models.User{})
		db.AutoMigrate(&models.User{})
		user := models.User{
			Email:    "admin@dev.com",
			Password: "admin1234",
			Role:     "Admin",
		}
		db.Create(&user)
	}

	var categories []models.Category
	err := db.Find(&categories).Error
	if len(categories) == 0 && err == nil {
		fmt.Println("Creating categories...")
		db.Migrator().DropTable(&models.Category{})
		db.AutoMigrate(&models.Category{})
		category := [...]string{"CPU", "GPU"}
		for i := 0; i < len(category); i++ {
			category := models.Category{
				Name: category[i],
			}

			categories = append(categories, category)
		}
		db.CreateInBatches(categories, len(category))
		fmt.Println("success")
	}

	numOfProducts := 1000
	products := make([]models.Product, numOfProducts)
	err = db.Find(&products).Limit(100).Error
	if len(products) == 0 && err == nil {
		fmt.Println("Creating products...")
		db.Migrator().DropTable(&models.Product{})
		db.AutoMigrate(&models.Product{})
		for i := 0; i < numOfProducts; i++ {
			product := models.Product{
				Name:       faker.Name(),
				Desc:       faker.Word(),
				Price:      rand.Intn(9999),
				Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
				CategoryID: uint(rand.Intn(2) + 1),
			}
			products = append(products, product)
		}
		db.CreateInBatches(products, 1000)
		fmt.Println("success")

	}
}