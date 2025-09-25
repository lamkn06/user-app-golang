package runtime

import "time"

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	// Enable WebSocket support
	Enabled bool `env:"WEBSOCKET_ENABLED" envDefault:"true"`
	
	// Maximum number of concurrent connections
	MaxConnections int `env:"WEBSOCKET_MAX_CONNECTIONS" envDefault:"1000"`
	
	// Read buffer size
	ReadBufferSize int `env:"WEBSOCKET_READ_BUFFER_SIZE" envDefault:"1024"`
	
	// Write buffer size
	WriteBufferSize int `env:"WEBSOCKET_WRITE_BUFFER_SIZE" envDefault:"1024"`
	
	// Message size limit
	MaxMessageSize int64 `env:"WEBSOCKET_MAX_MESSAGE_SIZE" envDefault:"512"`
	
	// Write wait timeout
	WriteWait time.Duration `env:"WEBSOCKET_WRITE_WAIT" envDefault:"10s"`
	
	// Pong wait timeout
	PongWait time.Duration `env:"WEBSOCKET_PONG_WAIT" envDefault:"60s"`
	
	// Ping period
	PingPeriod time.Duration `env:"WEBSOCKET_PING_PERIOD" envDefault:"54s"`
	
	// Allowed origins for CORS
	AllowedOrigins []string `env:"WEBSOCKET_ALLOWED_ORIGINS" envDefault:"*" envSeparator:","`
	
	// Enable authentication
	RequireAuth bool `env:"WEBSOCKET_REQUIRE_AUTH" envDefault:"false"`
	
	// Heartbeat interval
	HeartbeatInterval time.Duration `env:"WEBSOCKET_HEARTBEAT_INTERVAL" envDefault:"30s"`
	
	// Connection timeout
	ConnectionTimeout time.Duration `env:"WEBSOCKET_CONNECTION_TIMEOUT" envDefault:"60s"`
}
