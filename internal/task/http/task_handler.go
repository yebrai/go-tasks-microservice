package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yebrai/go-tasks-microservice/internal/task/creator"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
)

// TaskHandler maneja las peticiones HTTP relacionadas con tareas
type TaskHandler struct {
	commandBus cqrs.CommandBus
}

// NewTaskHandler crea una nueva instancia del handler
func NewTaskHandler(commandBus cqrs.CommandBus) *TaskHandler {
	return &TaskHandler{
		commandBus: commandBus,
	}
}

// CreateTaskRequest estructura de la petición
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date,omitempty"` // formato: "2006-01-02"
}

// CreateTaskResponse estructura de la respuesta
type CreateTaskResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// CreateTask maneja la creación de tareas
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest

	// Validar JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// Parsear fecha opcional
	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid due_date format",
				"message": "expected format: YYYY-MM-DD",
				"success": false,
			})
			return
		}
		dueDate = &parsed
	}

	// Crear comando
	cmd := creator.CreateTaskCommand{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDate,
	}

	// Ejecutar comando
	if err := h.commandBus.Dispatch(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create task",
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusCreated, CreateTaskResponse{
		Message: "Task created successfully",
		Success: true,
	})
}
