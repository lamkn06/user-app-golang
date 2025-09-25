package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/internal/service"
	"github.com/lamkn06/user-app-golang.git/pkg/logging"
)

// WebSocketRouter handles WebSocket connections
type WebSocketRouter struct {
	config    runtime.ServerConfig
	wsService service.WebSocketService
	logger    interface {
		Info(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
	}
}

// WebSocket upgrade configuration
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, implement proper origin checking
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NewWebSocketRouter creates a new WebSocket router
func NewWebSocketRouter(config runtime.ServerConfig, wsService service.WebSocketService) *WebSocketRouter {
	return &WebSocketRouter{
		config:    config,
		wsService: wsService,
		logger:    logging.NewSugaredLogger("websocket"),
	}
}

// Configure configures WebSocket routes
func (wr *WebSocketRouter) Configure(e *echo.Echo) {
	// WebSocket endpoint
	e.GET("/ws", wr.handleWebSocket)
	e.GET("/ws/:room", wr.handleWebSocketWithRoom)

	// WebSocket with authentication
	e.GET("/ws/auth", wr.handleAuthenticatedWebSocket)
	e.GET("/ws/auth/:room", wr.handleAuthenticatedWebSocketWithRoom)
}

// handleWebSocket handles WebSocket connection without authentication
func (wr *WebSocketRouter) handleWebSocket(c echo.Context) error {
	return wr.handleWebSocketConnection(c, "", "")
}

// handleWebSocketWithRoom handles WebSocket connection with room without authentication
func (wr *WebSocketRouter) handleWebSocketWithRoom(c echo.Context) error {
	room := c.Param("room")
	return wr.handleWebSocketConnection(c, room, "")
}

// handleAuthenticatedWebSocket handles WebSocket connection with authentication
func (wr *WebSocketRouter) handleAuthenticatedWebSocket(c echo.Context) error {
	// Extract user info from JWT token
	userID, username, err := wr.extractUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	return wr.handleWebSocketConnection(c, "", userID, username)
}

// handleAuthenticatedWebSocketWithRoom handles WebSocket connection with authentication and room
func (wr *WebSocketRouter) handleAuthenticatedWebSocketWithRoom(c echo.Context) error {
	// Extract user info from JWT token
	userID, username, err := wr.extractUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	room := c.Param("room")
	return wr.handleWebSocketConnection(c, room, userID, username)
}

// handleWebSocketConnection handles main WebSocket connection
func (wr *WebSocketRouter) handleWebSocketConnection(c echo.Context, room, userID string, username ...string) error {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		wr.logger.Error("Failed to upgrade connection:", err)
		return err
	}

	// Get username from parameters
	var userName string
	if len(username) > 0 {
		userName = username[0]
	}

	// Create client
	client := service.NewClient(conn, wr.wsService.(*service.Hub), userID, userName, room)

	// Register client with hub
	wr.wsService.RegisterClient(client)

	wr.logger.Info("New WebSocket client connected:", client.ID, "Room:", room, "User:", userName)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()

	return nil
}

// extractUserFromContext extracts user info from JWT token
func (wr *WebSocketRouter) extractUserFromContext(c echo.Context) (string, string, error) {
	// Get JWT token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
	}

	// Extract token from "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format")
	}

	token := tokenParts[1]
	fmt.Println("token", token)

	// TODO: Implement JWT token validation and user extraction
	// For now, return mock values
	// In real implementation, validate token and extract user info
	return "user123", "john_doe", nil
}

// WebSocketHandler handles different types of WebSocket messages
type WebSocketHandler struct {
	wsService service.WebSocketService
	logger    interface {
		Info(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
	}
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(wsService service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		wsService: wsService,
		logger:    logging.NewSugaredLogger("websocket-handler"),
	}
}

// HandleJoinRoom handles join room message
func (wh *WebSocketHandler) HandleJoinRoom(client *service.Client, room string) {
	// Remove from current room
	if client.Room != "" {
		wh.wsService.UnregisterClient(client)
	}

	// Add to new room
	client.Room = room
	wh.wsService.RegisterClient(client)

	// Broadcast to room
	wh.wsService.BroadcastToRoom(room, service.WebSocketMessage{
		Type:     service.MessageTypeJoin,
		Room:     room,
		UserID:   client.UserID,
		Username: client.Username,
		Content:  client.Username + " joined the room",
	})

	wh.logger.Info("User joined room:", client.Username, "Room:", room)
}

// HandleLeaveRoom handles leave room message
func (wh *WebSocketHandler) HandleLeaveRoom(client *service.Client) {
	if client.Room != "" {
		// Broadcast leave message
		wh.wsService.BroadcastToRoom(client.Room, service.WebSocketMessage{
			Type:     service.MessageTypeLeave,
			Room:     client.Room,
			UserID:   client.UserID,
			Username: client.Username,
			Content:  client.Username + " left the room",
		})

		wh.logger.Info("User left room:", client.Username, "Room:", client.Room)
	}

	// Unregister client
	wh.wsService.UnregisterClient(client)
}

// HandleSendMessage handles send message
func (wh *WebSocketHandler) HandleSendMessage(client *service.Client, content string) {
	if client.Room == "" {
		wh.logger.Error("Client not in any room:", client.ID)
		return
	}

	// Create message
	message := service.WebSocketMessage{
		Type:     service.MessageTypeMessage,
		Room:     client.Room,
		UserID:   client.UserID,
		Username: client.Username,
		Content:  content,
	}

	// Broadcast to room
	wh.wsService.BroadcastToRoom(client.Room, message)

	wh.logger.Debug("Message sent:", client.Username, "Room:", client.Room, "Content:", content)
}

// HandleGetUserList handles get user list
func (wh *WebSocketHandler) HandleGetUserList(client *service.Client) {
	var users []string
	if client.Room != "" {
		users = wh.wsService.GetRoomUsers(client.Room)
	} else {
		users = wh.wsService.GetAllUsers()
	}

	// Send user list to client
	message := service.WebSocketMessage{
		Type:   service.MessageTypeUserList,
		Room:   client.Room,
		UserID: client.UserID,
		Data:   users,
	}

	wh.wsService.SendToClient(client.ID, message)
}

// BroadcastSystemMessage sends system message
func (wh *WebSocketHandler) BroadcastSystemMessage(room string, content string) {
	message := service.WebSocketMessage{
		Type:     service.MessageTypeMessage,
		Room:     room,
		Username: "System",
		Content:  content,
	}

	if room != "" {
		wh.wsService.BroadcastToRoom(room, message)
	} else {
		wh.wsService.BroadcastToAll(message)
	}
}
