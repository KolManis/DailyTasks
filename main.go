package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

// func initDB() {
// 	var err error

// 	// Вариант 1 - URL формат
// 	connStr := "postgres://postgres:postgres@localhost:5433/dailytasks?sslmode=disable"

// 	// Вариант 2 - с явным указанием client_encoding
// 	// connStr := "user=postgres password=password123 dbname=dailytasks host=localhost port=5432 sslmode=disable client_encoding=UTF8"

// 	fmt.Println("Trying to connect with:", connStr)

// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Database connection failed:", err)
// 	}

// 	// Установи таймауты
// 	db.SetConnMaxLifetime(time.Minute * 3)
// 	db.SetMaxOpenConns(10)
// 	db.SetMaxIdleConns(10)
// 	db.SetConnMaxIdleTime(time.Minute * 1)

// 	// Дадим время на подключение
// 	time.Sleep(2 * time.Second)

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Database ping failed:", err)
// 	}

//		fmt.Println("✅ Connected to PostgreSQL!")
//	}
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

	fmt.Println("🔗 Connecting to PostgreSQL...")

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
