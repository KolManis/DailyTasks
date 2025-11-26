package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// router := gin.Default()

	// router.GET("/", func(c *gin.Context) {
	// 	c.String(200, "Hello World")
	// })

	// router.Run(":3000")

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
