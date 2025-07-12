// cmd/api/main.go - Simple Go API Application
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User represents a user in our API
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// In-memory storage (in production, use a real database)
var users = make(map[string]*User)

func main() {
	// Load configuration from environment variables
	port := getEnv("PORT", "9090")
	logLevel := getEnv("LOG_LEVEL", "info")
	appEnv := getEnv("APP_ENV", "development")

	// Initialize logging
	initLogging(logLevel)

	// Set Gin mode based on environment
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", healthHandler)

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/users", getUsersHandler)
		api.POST("/users", createUserHandler)
		api.GET("/users/:id", getUserHandler)
		api.PUT("/users/:id", updateUserHandler)
		api.DELETE("/users/:id", deleteUserHandler)
	}

	// Static files (if any)
	r.Static("/static", "./static")

	// Start server
	log.Printf("Starting Go API server on port %s (env: %s)", port, appEnv)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Health check handler
func healthHandler(c *gin.Context) {
	response := APIResponse{
		Status:  "ok",
		Message: "Go API is healthy",
		Data: gin.H{
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
			"uptime":    time.Since(startTime).String(),
		},
	}
	c.JSON(http.StatusOK, response)
}

// Get all users
func getUsersHandler(c *gin.Context) {
	var userList []*User
	for _, user := range users {
		userList = append(userList, user)
	}

	response := APIResponse{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    userList,
	}
	c.JSON(http.StatusOK, response)
}

// Create a new user
func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		response := APIResponse{
			Status:  "error",
			Message: "Invalid request body",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Generate ID and timestamp
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	// Store user
	users[user.ID] = &user

	// Save to persistent storage (example)
	saveUserToFile(&user)

	response := APIResponse{
		Status:  "success",
		Message: "User created successfully",
		Data:    user,
	}
	c.JSON(http.StatusCreated, response)
}

// Get a specific user
func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, exists := users[id]
	if !exists {
		response := APIResponse{
			Status:  "error",
			Message: "User not found",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := APIResponse{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	}
	c.JSON(http.StatusOK, response)
}

// Update a user
func updateUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, exists := users[id]
	if !exists {
		response := APIResponse{
			Status:  "error",
			Message: "User not found",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response := APIResponse{
			Status:  "error",
			Message: "Invalid request body",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Update fields
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	// Save to persistent storage
	saveUserToFile(user)

	response := APIResponse{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user,
	}
	c.JSON(http.StatusOK, response)
}

// Delete a user
func deleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	_, exists := users[id]
	if !exists {
		response := APIResponse{
			Status:  "error",
			Message: "User not found",
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	delete(users, id)

	response := APIResponse{
		Status:  "success",
		Message: "User deleted successfully",
	}
	c.JSON(http.StatusOK, response)
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Helper functions
var startTime = time.Now()

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initLogging(level string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[GO-API] ")
}

// Save user to persistent storage (example)
func saveUserToFile(user *User) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("/app/data", 0755); err != nil {
		log.Printf("Failed to create data directory: %v", err)
		return
	}

	// Save user to JSON file
	filename := fmt.Sprintf("/app/data/user_%s.json", user.ID)
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal user data: %v", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Failed to save user to file: %v", err)
		return
	}

	log.Printf("User saved to persistent storage: %s", filename)
}
