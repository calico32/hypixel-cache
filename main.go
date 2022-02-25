package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var mode string

func main() {
	loadEnv()

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.TrustedPlatform = gin.PlatformCloudflare
	r.SetTrustedProxies([]string{"172.0.0.1/16"})

	r.Use(cors)

	r.GET("/:type/:identifier", responseTimeStart, ensureAuthenticated, findCachedPlayer, fetchFromApi)
	r.Use(func(c *gin.Context) {
		c.JSON(404, NewFailure("not found"))
	})

	r.Run()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("API_SECRET") == "" {
		log.Fatal("$API_SECRET is empty or not set")
	}

	if os.Getenv("HYPIXEL_API_KEY") == "" {
		log.Fatal("$HYPIXEL_API_KEY is empty or not set")
	}
}
