package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	todos := []Todo{}
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	fmt.Println("-------------------------------")
	var x int = 5
	var p *int = &x
	fmt.Println(p)
	fmt.Println(*p)
	fmt.Println("-------------------------------")

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

		c.JSON(http.StatusNotFound, gin.H{"error": "Todo mot found"})
	})

	router.Run()
}
