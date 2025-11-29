package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var db *sql.DB

func initDB() {
	var err error

	// Загружаем .env файл
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Получаем настройки из .env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Формируем строку подключения из .env
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	fmt.Println("Connecting to PostgreSQL...")

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	migrationSQL, err := os.ReadFile("migrations/001_create_todos.sql")
	if err != nil {
		log.Fatal(" Failed to read migration:", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatal(" Failed to apply migration:", err)
	}

	fmt.Println(" Migration applied - todos table ready!")

	fmt.Println("Connected to PostgreSQL!")
}

func main() {
	initDB() // ← подключение к PostgreSQL
	defer db.Close()

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
		// Запрос к БД вместо памяти
		rows, err := db.Query("SELECT id, body, completed FROM todos ORDER BY id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		todos := []Todo{}

		for rows.Next() {
			var todo Todo
			err := rows.Scan(&todo.ID, &todo.Body, &todo.Completed)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Data parsing error"})
				return
			}
			todos = append(todos, todo)
		}

		c.JSON(http.StatusOK, todos)
	})

	router.POST("/api/todos", func(c *gin.Context) {
		todo := &Todo{}

		if err := c.BindJSON(todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		todo.Body = strings.TrimSpace(todo.Body)
		if todo.Body == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo body is required"})
			return
		}

		if len(todo.Body) > 1000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo body too long (max 1000 characters)"})
			return
		}

		err := db.QueryRow(
			"INSERT INTO todos (body, completed) VALUES ($1, $2) RETURNING id",
			todo.Body, false,
		).Scan(&todo.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}

		todo.Completed = false

		c.JSON(http.StatusCreated, todo)
	})

	router.PATCH("/api/todos/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, err = db.Exec("UPDATE todos SET completed =  NOT completed WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
			return
		}

		var updatedTodo Todo
		err := db.QueryRow(
			"SELECT id, body, completed FROM todos WHERE id = $1",
			id).Scan(&updatedTodo.ID, &updatedTodo.Body, &updatedTodo.Completed)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(http.StatusOK, updatedTodo)
	})

	router.DELETE("/api/todos/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		result, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
