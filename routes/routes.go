package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/app/v2/cache"
	"github.com/sing3demons/app/v2/controllers"
	"github.com/sing3demons/app/v2/database"
	"github.com/sing3demons/app/v2/middleware"
	"github.com/sing3demons/app/v2/store"
)

//Serve - middleware
func Serve(r *gin.Engine) {

	db := database.GetDB()
	cacher := cache.NewCacher(cache.NewCacherConfig())
	v1 := r.Group("/api/v1")

	store := store.NewGormStorm(db)

	authenticate := middleware.JwtVerify()
	authorize := middleware.Authorize()

	authGroup := v1.Group("/auth")
	authController := controllers.NewAuthHandler(store)
	{
		authGroup.POST("/register", authController.SignUp)
		authGroup.POST("/login", middleware.Login)
		authGroup.Use(authenticate)
		authGroup.GET("/profile", authController.GetProfile)
		authGroup.PATCH("/profile/:id", authController.UpdateImageProfile)
		authGroup.PUT("/profile/:id", authController.UpdateProfile)
	}

	categoryGroup := v1.Group("/categories")
	categoryController := controllers.NewCategoryHandler(store)
	{
		categoryGroup.GET("", categoryController.FindAll)
		categoryGroup.GET("/:id", categoryController.FindOne)
		categoryGroup.POST("", categoryController.Create)
		categoryGroup.PUT("/:id", categoryController.Update)
		categoryGroup.DELETE("/:id", categoryController.Delete)
	}

	productGroup := v1.Group("/products")
	productController := controllers.NewProductHandler(store, cacher)
	productGroup.GET("", productController.FindAll)
	productGroup.GET("/:id", productController.FindOne)

	productGroup.Use(authenticate, authorize)
	{
		productGroup.POST("", productController.Create)
		productGroup.PUT("/:id", productController.UpdateAll)
		productGroup.DELETE("/:id", productController.Delete)
	}

	usersGroup := v1.Group("/users")
	usersController := controllers.Users{DB: db}
	usersGroup.Use(authenticate, authorize)
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
