package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sing3demons/app/v2/database"
	"github.com/sing3demons/app/v2/routes"
	"github.com/sing3demons/app/v2/seeds"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	docs "github.com/sing3demons/app/v2/docs"
	swagger "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	buildCommit = "dev"
	buildTime   = time.Now().String()
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
}

// @title Swagger GO-PRODUCT API
// @version 1.0
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	livenessProbe()

	database.InitDB()
	seeds.Load()

	r := routes.NewMyRouter()
	r.Static("/uploads", "./uploads")

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagger.Handler))

	r.GET("/healthz", health)

	r.GET("/x", buildX)

	//สร้าง folder
	uploadDirs := [...]string{"products", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(r)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swagger.Handler))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		fmt.Printf("listen and serve on http://localhost:8080/swagger/index.html\n")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

// @Accept  json
// @Produce  json
// @Success 200
// @Router /healthz [get]
func health(c *gin.Context) {
	c.Status(200)
}

// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]any
// @Router /x [get]
func buildX(c *gin.Context) {
	c.JSON(200, gin.H{
		"build_commit": buildCommit,
		"build_time":   buildTime,
	})
}

func livenessProbe() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")
}
