package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/pkg/events"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		log.Printf("🔍 WebSocket origin check: %s", r.Header.Get("Origin"))
		return true
	},
}

// WebSocketHandler manages WebSocket connections and event broadcasting
type WebSocketHandler struct {
	clients    map[*websocket.Conn]bool
	clientsMux sync.RWMutex
	eventBus   events.EventBus
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(eventBus events.EventBus) *WebSocketHandler {
	return &WebSocketHandler{
		clients:  make(map[*websocket.Conn]bool),
		eventBus: eventBus,
	}
}

// EventMessage represents a message sent to WebSocket clients
type EventMessage struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	AggregateID string      `json:"aggregateId"`
	Timestamp   time.Time   `json:"timestamp"`
	Payload     interface{} `json:"payload"`
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	log.Printf("🔗 WebSocket connection attempt from: %s", c.ClientIP())

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add client to the list
	h.clientsMux.Lock()
	h.clients[conn] = true
	clientCount := len(h.clients)
	h.clientsMux.Unlock()

	log.Printf("✅ WebSocket client connected. Total clients: %d", clientCount)

	// Remove client when disconnected
	defer func() {
		h.clientsMux.Lock()
		delete(h.clients, conn)
		h.clientsMux.Unlock()
		log.Printf("WebSocket client disconnected. Total clients: %d", len(h.clients))
	}()

	// Send welcome message
	welcomeMsg := EventMessage{
		ID:        "welcome",
		Type:      "connection.established",
		Timestamp: time.Now(),
		Payload:   map[string]string{"message": "Connected to task event stream"},
	}
	conn.WriteJSON(welcomeMsg)

	// Keep connection alive and handle client messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// BroadcastEvent sends an event to all connected WebSocket clients
func (h *WebSocketHandler) BroadcastEvent(event task.DomainEvent) {
	if len(h.clients) == 0 {
		return // No clients connected
	}

	message := EventMessage{
		ID:          fmt.Sprintf("event_%d", time.Now().UnixNano()),
		Type:        event.EventName(),
		AggregateID: event.AggregateID(),
		Timestamp:   event.OccurredOn(),
		Payload:     h.extractEventPayload(event),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal event message: %v", err)
		return
	}

	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	var disconnectedClients []*websocket.Conn

	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			log.Printf("Failed to send message to WebSocket client: %v", err)
			disconnectedClients = append(disconnectedClients, client)
		}
	}

	// Clean up disconnected clients
	for _, client := range disconnectedClients {
		delete(h.clients, client)
		client.Close()
	}

	log.Printf("Broadcasted event %s to %d clients", event.EventName(), len(h.clients))
}

// extractEventPayload extracts the relevant payload from different event types
func (h *WebSocketHandler) extractEventPayload(event task.DomainEvent) interface{} {
	switch e := event.(type) {
	case *task.TaskCreatedEvent:
		return map[string]interface{}{
			"task_id":     e.TaskID,
			"title":       e.Title,
			"description": e.Description,
			"due_date":    e.DueDate,
		}
	case *task.TaskCompletedEvent:
		return map[string]interface{}{
			"task_id": e.TaskID,
		}
	case *task.TaskCancelledEvent:
		return map[string]interface{}{
			"task_id": e.TaskID,
		}
	default:
		return map[string]interface{}{
			"aggregate_id": event.AggregateID(),
		}
	}
}

// StartEventListener starts listening for events from RabbitMQ (if implemented)
func (h *WebSocketHandler) StartEventListener(ctx context.Context) {
	// This would connect to RabbitMQ to listen for events
	// For now, we'll rely on the command handlers to call BroadcastEvent directly
	log.Println("WebSocket event listener started (direct mode)")
}
