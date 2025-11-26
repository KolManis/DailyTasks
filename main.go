package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	todos := []Todo{}
	router := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080" // порт по умолчанию если в .env нет PORT
	// }

	router.GET("/api/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, todos)
	})

	router.POST("/api/todos", func(c *gin.Context) {
		todo := &Todo{}

		if err := c.BindJSON(todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		if todo.Body == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo body is required"})
			return
		}

		if len(todos) > 0 {
			todo.ID = todos[len(todos)-1].ID + 1
		} else {
			todo.ID = 1
		}

		todos = append(todos, *todo)

		c.JSON(http.StatusCreated, todo)
	})

	router.PATCH("/api/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				c.JSON(http.StatusOK, todos[i])
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	})

	router.DELETE("/api/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"success": "true"})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	})

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
