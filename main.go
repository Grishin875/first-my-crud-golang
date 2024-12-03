package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var sliceTask []Task = []Task{}

// create task
func createNewTask(c *gin.Context) {
	var task Task
	err := c.BindJSON(&task)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = len(sliceTask) + 1
	sliceTask = append(sliceTask, task)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created",
		"task":    task,
	})
}

// delete task
func deleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for index, task := range sliceTask {
		if task.ID == id {
			sliceTask = append(sliceTask[:index], sliceTask[index+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Task deleted",
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Task not found",
	})
}

// update task
func putTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTask Task
	err = c.BindJSON(&updatedTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range sliceTask {
		if task.ID == id {
			sliceTask[i].Status = updatedTask.Status // Обновляем только статус
			c.JSON(http.StatusOK, gin.H{
				"message": "Task updated",
				"task":    sliceTask[i],
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

// get tasks
func getTasks(c *gin.Context) {
	status := c.DefaultQuery("status", "") // Получаем параметр status (если есть)

	var filteredTasks []Task
	for _, task := range sliceTask {
		if status == "" || task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	if len(filteredTasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": filteredTasks,
	})
}

func main() {
	// Создаем новый роутер
	router := gin.Default()

	// Настройка CORS с явными параметрами
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                            // Разрешаем все домены
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Разрешаем методы
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Роуты API
	router.GET("/tasks", getTasks)
	router.POST("/tasks", createNewTask)
	router.PUT("/tasks/:id", putTask)
	router.DELETE("/tasks/:id", deleteTask)

	// Запуск сервера на порту 8080
	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
