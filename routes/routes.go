package routes

import (
	"app/config"
	"app/controllers"
	"app/middleware"

	"github.com/gin-gonic/gin"
)

//Serve - middleware
func Serve(r *gin.Engine) {

	db := config.GetDB()
	v1 := r.Group("/api/v1")

	authenticate := middleware.JwtVerify()

	authGroup := v1.Group("/auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.GET("/profile", authenticate, authController.GetProfile)
		authGroup.POST("register", authController.SignUp)
		authGroup.POST("/login", middleware.Login)
		authGroup.PATCH("/profile/:id", authenticate, authController.UpdateImageProfile)
		authGroup.PUT("/profile/:id", authenticate, authController.UpdateProfile)
	}

	categoryGroup := v1.Group("/categories")
	categoryController := controllers.Category{DB: db}
	{
		categoryGroup.GET("", categoryController.FindAll)
		categoryGroup.GET("/:id", categoryController.FindOne)
		categoryGroup.POST("", categoryController.Create)
		categoryGroup.PUT("/:id", categoryController.Update)
		categoryGroup.DELETE("/:id", categoryController.Delete)
	}

	productGroup := v1.Group("/products")
	productController := controllers.Product{DB: db}
	{
		productGroup.GET("", productController.FindAll)
		productGroup.GET("/:id", productController.FindOne)
		productGroup.POST("", productController.Create)
		productGroup.PUT("/:id", productController.UpdateAll)
		productGroup.DELETE("/:id", productController.Delete)
	}

	usersGroup := v1.Group("/users")
	usersController := controllers.Users{DB: db}
	usersGroup.Use(authenticate)
	{
		usersGroup.GET("", usersController.FindAll)
		usersGroup.POST("", usersController.Create)
		usersGroup.GET("/:id", usersController.FindOne)
		usersGroup.PUT("/:id", usersController.Update)
		usersGroup.DELETE("/:id", usersController.Delete)
		usersGroup.PATCH("/:id/promote", usersController.Promote)
		usersGroup.PATCH("/:id/demote", usersController.Demote)
	}
}
