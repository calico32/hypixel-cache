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
	r.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.GET("/:type/:identifier", ensureAuthenticated, findCachedPlayer, fetchFromApi)

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
