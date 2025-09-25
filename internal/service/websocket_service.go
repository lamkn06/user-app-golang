package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// MessageType defines message types for WebSocket
type MessageType string

const (
	MessageTypeJoin     MessageType = "join"
	MessageTypeLeave    MessageType = "leave"
	MessageTypeMessage  MessageType = "message"
	MessageTypePing     MessageType = "ping"
	MessageTypePong     MessageType = "pong"
	MessageTypeError    MessageType = "error"
	MessageTypeUserList MessageType = "user_list"
)

// WebSocketMessage defines message structure
type WebSocketMessage struct {
	Type      MessageType `json:"type"`
	Room      string      `json:"room,omitempty"`
	Content   string      `json:"content,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	Username  string      `json:"username,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// Client represents a WebSocket connection
type Client struct {
	ID       string
	UserID   string
	Username string
	Room     string
	Conn     *websocket.Conn
	Send     chan WebSocketMessage
	Hub      *Hub
}

// Hub manages all WebSocket connections
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan WebSocketMessage

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Room management
	rooms map[string]map[*Client]bool

	// Mutex for concurrent access
	mutex sync.RWMutex

	// Logger
	logger interface {
		Info(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
	}
}

// WebSocketService interface
type WebSocketService interface {
	Run(ctx context.Context)
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastToRoom(room string, message WebSocketMessage)
	BroadcastToAll(message WebSocketMessage)
	SendToClient(clientID string, message WebSocketMessage) error
	GetRoomUsers(room string) []string
	GetAllUsers() []string
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan WebSocketMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

// Run starts the Hub to handle WebSocket connections
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("WebSocket Hub shutting down...")
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.handleMessage(message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	// Add to room if specified
	if client.Room != "" {
		if h.rooms[client.Room] == nil {
			h.rooms[client.Room] = make(map[*Client]bool)
		}
		h.rooms[client.Room][client] = true
	}

	h.logger.Info("Client registered:", client.ID, "in room:", client.Room)

	// Send welcome message to new client
	welcomeMsg := WebSocketMessage{
		Type:      MessageTypeMessage,
		UserID:    client.UserID,
		Username:  "System",
		Content:   "Welcome to WebSocket! Your connection ID: " + client.ID,
		Timestamp: getCurrentTimestamp(),
	}
	select {
	case client.Send <- welcomeMsg:
	default:
		// Client channel is full, skip
	}

	// Notify room about new user
	if client.Room != "" {
		// Send join notification to room
		h.broadcastToRoomUnsafe(client.Room, WebSocketMessage{
			Type:      MessageTypeJoin,
			Room:      client.Room,
			UserID:    client.UserID,
			Username:  client.Username,
			Content:   client.Username + " joined the room",
			Timestamp: getCurrentTimestamp(),
		})

		// Send updated user list to room
		h.broadcastToRoomUnsafe(client.Room, WebSocketMessage{
			Type:      MessageTypeUserList,
			Room:      client.Room,
			UserID:    client.UserID,
			Username:  client.Username,
			Timestamp: getCurrentTimestamp(),
			Data:      h.getRoomUserListUnsafe(client.Room),
		})
	}
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// Remove from room
		if client.Room != "" && h.rooms[client.Room] != nil {
			delete(h.rooms[client.Room], client)
			if len(h.rooms[client.Room]) == 0 {
				delete(h.rooms, client.Room)
			}
		}

		h.logger.Info("Client unregistered:", client.ID)

		// Notify room about user leaving
		if client.Room != "" {
			// Send leave notification to room
			h.broadcastToRoomUnsafe(client.Room, WebSocketMessage{
				Type:      MessageTypeLeave,
				Room:      client.Room,
				UserID:    client.UserID,
				Username:  client.Username,
				Content:   client.Username + " left the room",
				Timestamp: getCurrentTimestamp(),
			})

			// Send updated user list to room
			h.broadcastToRoomUnsafe(client.Room, WebSocketMessage{
				Type:      MessageTypeUserList,
				Room:      client.Room,
				UserID:    client.UserID,
				Username:  client.Username,
				Timestamp: getCurrentTimestamp(),
				Data:      h.getRoomUserListUnsafe(client.Room),
			})
		}
	}
}

// handleMessage handles messages from clients
func (h *Hub) handleMessage(message WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Add timestamp to message
	message.Timestamp = getCurrentTimestamp()

	switch message.Type {
	case MessageTypeJoin:
		// Handle join room logic
		h.logger.Debug("Join message received:", message)
		if message.Room != "" {
			// Find client and update room
			for client := range h.clients {
				if client.UserID == message.UserID {
					// Remove from current room if any
					if client.Room != "" && h.rooms[client.Room] != nil {
						delete(h.rooms[client.Room], client)
					}
					// Add to new room
					client.Room = message.Room
					if h.rooms[message.Room] == nil {
						h.rooms[message.Room] = make(map[*Client]bool)
					}
					h.rooms[message.Room][client] = true

					// Notify room about join
					h.broadcastToRoomUnsafe(message.Room, WebSocketMessage{
						Type:      MessageTypeJoin,
						Room:      message.Room,
						UserID:    client.UserID,
						Username:  client.Username,
						Content:   client.Username + " joined the room",
						Timestamp: message.Timestamp,
					})

					// Send user list
					h.broadcastToRoomUnsafe(message.Room, WebSocketMessage{
						Type:      MessageTypeUserList,
						Room:      message.Room,
						Timestamp: message.Timestamp,
						Data:      h.getRoomUserListUnsafe(message.Room),
					})
					break
				}
			}
		}
	case MessageTypeLeave:
		// Handle leave room logic
		h.logger.Debug("Leave message received:", message)
		for client := range h.clients {
			if client.UserID == message.UserID && client.Room != "" {
				room := client.Room
				// Remove from room
				if h.rooms[room] != nil {
					delete(h.rooms[room], client)
				}
				client.Room = ""

				// Notify room about leave
				h.broadcastToRoomUnsafe(room, WebSocketMessage{
					Type:      MessageTypeLeave,
					Room:      room,
					UserID:    client.UserID,
					Username:  client.Username,
					Content:   client.Username + " left the room",
					Timestamp: message.Timestamp,
				})

				// Send updated user list
				h.broadcastToRoomUnsafe(room, WebSocketMessage{
					Type:      MessageTypeUserList,
					Room:      room,
					Timestamp: message.Timestamp,
					Data:      h.getRoomUserListUnsafe(room),
				})
				break
			}
		}
	case MessageTypeMessage:
		// Broadcast message to room
		h.logger.Debug("Message received:", message)
		if message.Room != "" {
			h.broadcastToRoomUnsafe(message.Room, message)
		} else {
			h.broadcastToAllUnsafe(message)
		}
	case MessageTypePing:
		// Handle ping - send pong back
		h.logger.Debug("Ping received from:", message.UserID)
		for client := range h.clients {
			if client.UserID == message.UserID {
				pongMsg := WebSocketMessage{
					Type:      MessageTypePong,
					UserID:    client.UserID,
					Username:  client.Username,
					Timestamp: message.Timestamp,
				}
				select {
				case client.Send <- pongMsg:
				default:
					// Client channel is full, skip
				}
				break
			}
		}
	default:
		h.logger.Debug("Unknown message type:", message.Type)
		// Send error response
		for client := range h.clients {
			if client.UserID == message.UserID {
				errorMsg := WebSocketMessage{
					Type:      MessageTypeError,
					UserID:    client.UserID,
					Username:  client.Username,
					Content:   "Unknown message type: " + string(message.Type),
					Timestamp: message.Timestamp,
				}
				select {
				case client.Send <- errorMsg:
				default:
					// Client channel is full, skip
				}
				break
			}
		}
	}
}

// BroadcastToRoom sends message to all clients in room
func (h *Hub) BroadcastToRoom(room string, message WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	h.broadcastToRoomUnsafe(room, message)
}

// broadcastToRoomUnsafe sends message to room (without lock)
func (h *Hub) broadcastToRoomUnsafe(room string, message WebSocketMessage) {
	if clients, ok := h.rooms[room]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				delete(h.clients, client)
				delete(h.rooms[room], client)
				close(client.Send)
			}
		}
	}
}

// BroadcastToAll sends message to all clients
func (h *Hub) BroadcastToAll(message WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	h.broadcastToAllUnsafe(message)
}

// broadcastToAllUnsafe sends message to all clients (without lock)
func (h *Hub) broadcastToAllUnsafe(message WebSocketMessage) {
	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			delete(h.clients, client)
			if client.Room != "" && h.rooms[client.Room] != nil {
				delete(h.rooms[client.Room], client)
			}
			close(client.Send)
		}
	}
}

// SendToClient sends message to specific client
func (h *Hub) SendToClient(clientID string, message WebSocketMessage) error {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.ID == clientID {
			select {
			case client.Send <- message:
				return nil
			default:
				return websocket.ErrCloseSent
			}
		}
	}
	return websocket.ErrBadHandshake
}

// GetRoomUsers gets list of users in room
func (h *Hub) GetRoomUsers(room string) []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.getRoomUserListUnsafe(room)
}

// getRoomUserListUnsafe gets list of users in room (without lock)
func (h *Hub) getRoomUserListUnsafe(room string) []string {
	var users []string
	if clients, ok := h.rooms[room]; ok {
		for client := range clients {
			if client.Username != "" {
				users = append(users, client.Username)
			}
		}
	}
	return users
}

// GetAllUsers gets list of all users
func (h *Hub) GetAllUsers() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var users []string
	for client := range h.clients {
		if client.Username != "" {
			users = append(users, client.Username)
		}
	}
	return users
}

// RegisterClient registers a client
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// readPump handles reading messages from WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	// Set read limit and pong handler
	c.Conn.SetReadLimit(512)
	c.Conn.SetPongHandler(func(string) error {
		return nil
	})

	for {
		var message WebSocketMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Set client info in message
		message.UserID = c.UserID
		message.Username = c.Username

		// Send to hub for processing
		select {
		case c.Hub.broadcast <- message:
		default:
			// Hub channel is full, skip message
		}
	}
}

// writePump handles sending messages to WebSocket connection
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		c.Conn.SetWriteDeadline(getWriteDeadline())
		if err := c.Conn.WriteJSON(message); err != nil {
			return
		}
	}
}

// Helper functions
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func getWriteDeadline() time.Time {
	return time.Now().Add(10 * time.Second)
}

// NewClient creates a new client
func NewClient(conn *websocket.Conn, hub *Hub, userID, username, room string) *Client {
	return &Client{
		ID:       uuid.New().String(),
		UserID:   userID,
		Username: username,
		Room:     room,
		Conn:     conn,
		Send:     make(chan WebSocketMessage, 256),
		Hub:      hub,
	}
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService() WebSocketService {
	return NewHub()
}
