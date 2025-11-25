package main

import (
    "github.com/gin-gonic/gin"
    "go-tasks-api/internal/auth"
    "go-tasks-api/internal/handlers"
)

func main() {
	//new gin router with default middleware (logging + recovery)
    r := gin.Default()

    // Public route: login
    r.POST("/login", func(c *gin.Context) {
        // For practice, hardcode a user
        token, _ := auth.GenerateToken("user123")
        c.JSON(200, gin.H{"token": token})
    })

    // Protected routes (group with base path /api/v1)
    api := r.Group("/api/v1")
	// Applies your JWTMiddleware() to all routes
    api.Use(auth.JWTMiddleware())
    {
        api.POST("/tasks", handlers.CreateTask)
        api.GET("/tasks", handlers.GetTasks)
        // Add PUT, PATCH, DELETE similarly
    }

    r.Run(":8080")
}
