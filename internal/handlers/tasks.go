package handlers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-tasks-api/internal/models"
)
// key: string -> value: models.Task structs
var tasks = map[string]models.Task{}

//HTTP handler that receives the Gin context c
func CreateTask(c *gin.Context){
	var input models.Task
	//parse request body as JSON into input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
    input.ID = id
    input.CreatedAt = time.Now()
    input.UpdatedAt = time.Now()
    input.UserID = c.GetString("userId") // from JWT middleware
    tasks[id] = input
    c.JSON(http.StatusCreated, input)
}

//GET
func GetTasks(c *gin.Context){
	userId := c.GetString("userId")
	//slice of tasks
    userTasks := []models.Task{}
	// _ discard the key
    for _, t := range tasks {
        if t.UserID == userId {
            userTasks = append(userTasks, t)
        }
    }
    c.JSON(http.StatusOK, userTasks)
}

//PUT
func UpdateTask(c *gin.Context){
	id := c.Param("id")
    userId := c.GetString("userId")

    task, exists := tasks[id]
    if !exists || task.UserID != userId {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    var input models.Task
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Replace fields
    input.ID = id
    input.UserID = userId
    input.CreatedAt = task.CreatedAt
    input.UpdatedAt = time.Now()

    tasks[id] = input
    c.JSON(http.StatusOK, input)
}

func PatchTask(c *gin.Context){
	id := c.Param("id")
    userId := c.GetString("userId")

    task, exists := tasks[id]
    if !exists || task.UserID != userId {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    var input map[string]interface{}
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Update only provided fields
    if title, ok := input["title"].(string); ok {
        task.Title = title
    }
    if desc, ok := input["description"].(string); ok {
        task.Description = desc
    }
    if status, ok := input["status"].(string); ok {
        task.Status = status
    }

    task.UpdatedAt = time.Now()
    tasks[id] = task
    c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
    id := c.Param("id")
    userId := c.GetString("userId")

    task, exists := tasks[id]
    if !exists || task.UserID != userId {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    delete(tasks, id)
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

