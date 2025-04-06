package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"net/http"
)

// Todo struct represents the Todo model in the database
type Todo struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Initialize the database connection
var db *gorm.DB
var err error

// Connect to PostgreSQL
func initDB() {
	// Set your PostgreSQL connection details here
	dsn := "user=postgres password=123456 dbname=todo_db port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect to the database:", err)
		panic("failed to connect to the database")
	}

	// Auto-migrate the Todo model
	db.AutoMigrate(&Todo{})
}

// Get all todos
func getTodos(c *gin.Context) {
	var todos []Todo
	if err := db.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// Add a new todo
func addTodo(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	if err := db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add todo"})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

// Update a todo's status (mark it as done)
func markDone(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	todo.Done = true
	if err := db.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

// Delete a todo
func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Todo{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func main() {
	// Initialize the database
	initDB()

	// Create a Gin router
	r := gin.Default()

	// Serve static files (HTML, CSS, JS)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Set up routes for API
	r.GET("/api/todos", getTodos)
	r.POST("/api/todos", addTodo)
	r.PUT("/api/todos/:id", markDone)
	r.DELETE("/api/todos/:id", deleteTodo)

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Start the server on port 8080
	r.Run(":8081")
}
