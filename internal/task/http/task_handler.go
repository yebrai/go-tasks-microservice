package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/internal/task/creator"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
)

// TaskHandler maneja las peticiones HTTP relacionadas con tareas
type TaskHandler struct {
	commandBus cqrs.CommandBus
	repository task.Repository
	wsHandler  *WebSocketHandler
}

// NewTaskHandler crea una nueva instancia del handler
func NewTaskHandler(commandBus cqrs.CommandBus, repository task.Repository, wsHandler *WebSocketHandler) *TaskHandler {
	return &TaskHandler{
		commandBus: commandBus,
		repository: repository,
		wsHandler:  wsHandler,
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

// GetTasks maneja la consulta de todas las tareas
func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.repository.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch tasks",
			"message": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    tasks,
		"success": true,
	})
}

// GetTask maneja la consulta de una tarea específica
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "task id is required",
			"success": false,
		})
		return
	}

	taskID, err := task.NewID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid task id",
			"message": err.Error(),
			"success": false,
		})
		return
	}

	foundTask, err := h.repository.FindByID(c.Request.Context(), string(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "task not found",
			"message": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    foundTask,
		"success": true,
	})
}

// UpdateTaskRequest estructura de la petición de actualización
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

// UpdateTask maneja la actualización de tareas
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "task id is required",
			"success": false,
		})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// Por ahora, solo soportamos cambio de status a completed
	if req.Status != nil && *req.Status == "completed" {
		cmd := creator.CompleteTaskCommand{ID: id}
		if err := h.commandBus.Dispatch(c.Request.Context(), cmd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to complete task",
				"message": err.Error(),
				"success": false,
			})
			return
		}

		// Enviar evento WebSocket para tarea completada
		if h.wsHandler != nil {
			updatedTask, err := h.repository.FindByID(c.Request.Context(), id)
			if err == nil {
				event := task.NewTaskCompletedEvent(updatedTask)
				h.wsHandler.BroadcastEvent(event)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"success": true,
	})
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

	// Encontrar la tarea recién creada para enviar evento WebSocket
	// Como no tenemos el ID aquí, buscaremos la tarea más reciente con el mismo título
	if h.wsHandler != nil {
		log.Printf("🔍 Attempting to broadcast task creation event")
		tasks, err := h.repository.FindAll(c.Request.Context())
		if err == nil && len(tasks) > 0 {
			// Buscar la tarea más reciente con el mismo título
			for i := len(tasks) - 1; i >= 0; i-- {
				if tasks[i].Title == req.Title {
					event := task.NewTaskCreatedEvent(tasks[i])
					log.Printf("📡 Broadcasting task created event for task: %s", tasks[i].ID)
					h.wsHandler.BroadcastEvent(event)
					break
				}
			}
		} else {
			log.Printf("❌ Failed to fetch tasks for WebSocket broadcast: %v", err)
		}
	} else {
		log.Printf("⚠️  WebSocket handler is nil, cannot broadcast event")
	}

	// Respuesta exitosa
	c.JSON(http.StatusCreated, CreateTaskResponse{
		Message: "Task created successfully",
		Success: true,
	})
}
