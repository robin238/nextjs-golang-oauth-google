package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	InitDB()

	r := gin.Default()

	r.GET("/auth/google", func(c *gin.Context) {
		GoogleLogin(c.Writer, c.Request)
	})

	r.GET("/auth/google/callback", func(c *gin.Context) {
		GoogleCallback(c.Writer, c.Request)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome dashboard"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback default
	}

	r.Run(":" + port)
}
