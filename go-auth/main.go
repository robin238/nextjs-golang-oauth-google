package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
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

	r.Run(":" + os.Getenv("PORT"))
}
