package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	r.GET("/auth/google", func(c *gin.Context) {
		GoogleLogin(c.Writer, c.Request)
	})

	r.GET("/auth/google/callback", func(c *gin.Context) {
		GoogleCallback(c.Writer, c.Request)
	})

	r.POST("/logout", AuthMiddleware(), func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		RevokeToken(tokenString)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	})

	r.GET("/dashboard", AuthMiddleware(), func(c *gin.Context) {
		email := c.GetString("email")
		name := c.GetString("name")
		role := c.GetString("role")

		if role == "admin" {
			c.JSON(200, gin.H{"message": "Welcome Admin", "email": email, "name": name, "role": role})
			return
		}

		c.JSON(200, gin.H{"message": "Welcome User", "email": email, "name": name, "role": role})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback default
	}

	r.Run(":" + port)
}
