package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
	"github.com/yebrai/go-tasks-microservice/pkg/events"
)

// Server maneja el servidor HTTP
type Server struct {
	commandBus cqrs.CommandBus
	repository task.Repository
	handler    *TaskHandler
	wsHandler  *WebSocketHandler
}

// NewServer crea una nueva instancia del servidor HTTP
func NewServer(commandBus cqrs.CommandBus, repository task.Repository, eventBus events.EventBus) *Server {
	wsHandler := NewWebSocketHandler(eventBus)
	return &Server{
		commandBus: commandBus,
		repository: repository,
		handler:    NewTaskHandler(commandBus, repository, wsHandler),
		wsHandler:  wsHandler,
	}
}

// Handler retorna el handler HTTP configurado
func (s *Server) Handler() http.Handler {
	// Configurar Gin en modo release
	gin.SetMode(gin.ReleaseMode)

	// Crear router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(s.loggingMiddleware())
	router.Use(s.corsMiddleware())

	// Configurar rutas
	s.setupRoutes(router)

	return router
}

// setupRoutes configura todas las rutas del API
func (s *Server) setupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", s.healthCheck)

	// WebSocket endpoint
	router.GET("/ws/events", s.wsHandler.HandleWebSocket)

	// API v1
	api := router.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", s.handler.GetTasks)
			tasks.POST("", s.handler.CreateTask)
			tasks.GET("/:id", s.handler.GetTask)
			tasks.PUT("/:id", s.handler.UpdateTask)
		}
	}
}

// healthCheck endpoint de salud
func (s *Server) healthCheck(c *gin.Context) {
	// Verificar conectividad de la base de datos intentando obtener tareas
	dbStatus := "connected"
	_, err := s.repository.FindAll(c.Request.Context())
	if err != nil {
		dbStatus = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"service":  "go-tasks-microservice",
		"database": dbStatus,
		"rabbitmq": "connected", // Por ahora asumimos que está conectado, mejorar en el futuro
	})
}

// loggingMiddleware middleware de logging
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s - %d (%v)\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	})
}

// corsMiddleware middleware CORS
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
