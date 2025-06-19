package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
)

// Server maneja el servidor HTTP
type Server struct {
	commandBus cqrs.CommandBus
	handler    *TaskHandler
}

// NewServer crea una nueva instancia del servidor HTTP
func NewServer(commandBus cqrs.CommandBus) *Server {
	return &Server{
		commandBus: commandBus,
		handler:    NewTaskHandler(commandBus),
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

	// API v1
	api := router.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("", s.handler.CreateTask)
		}
	}
}

// healthCheck endpoint de salud
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "go-tasks-microservice",
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
