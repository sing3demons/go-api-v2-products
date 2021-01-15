package routes

import (
	"app/config"
	"app/controllers"

	"github.com/gin-gonic/gin"
)

//Serve - middleware
func Serve(r *gin.Engine) {
	db := config.GetDB()
	productGroup := r.Group("/api/v1/products")
	productController := controllers.Product{DB: db}
	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne)
		productGroup.POST("", productController.Create)
	}
}
