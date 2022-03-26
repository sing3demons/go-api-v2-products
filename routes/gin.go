package routes

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewMyRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}
	r.Use(cors.New(config))
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return r
}
