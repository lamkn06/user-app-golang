package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/pkg/logging"
	"go.uber.org/zap"
)

// WebSocketAuthMiddleware middleware for authenticating WebSocket connections
type WebSocketAuthMiddleware struct {
	jwtConfig runtime.JWTConfig
	logger    interface {
		Info(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
	}
}

// NewWebSocketAuthMiddleware creates a new WebSocket auth middleware
func NewWebSocketAuthMiddleware(jwtConfig runtime.JWTConfig) *WebSocketAuthMiddleware {
	return &WebSocketAuthMiddleware{
		jwtConfig: jwtConfig,
		logger:    logging.NewSugaredLogger("websocket-auth"),
	}
}

// AuthenticateWebSocket authenticates WebSocket connection
func (wm *WebSocketAuthMiddleware) AuthenticateWebSocket(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract token from query parameter or header
		token := wm.extractToken(c)
		if token == "" {
			wm.logger.Error("No token provided for WebSocket connection")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authentication required",
			})
		}

		// Validate token
		claims, err := wm.validateToken(token)
		if err != nil {
			wm.logger.Error("Token validation failed:", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
		}

		// Extract user info from claims
		userID, username := wm.extractUserInfo(claims)
		if userID == "" {
			wm.logger.Error("No user ID found in token")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid user information",
			})
		}

		// Store user info in context
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("token_claims", claims)

		wm.logger.Debug("WebSocket authentication successful for user:", username)
		return next(c)
	}
}

// extractToken extracts token from request
func (wm *WebSocketAuthMiddleware) extractToken(c echo.Context) string {
	// Try to get token from Authorization header first
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try to get token from query parameter (common for WebSocket)
	token := c.QueryParam("token")
	if token != "" {
		return token
	}

	// Try to get token from cookie
	cookie, err := c.Cookie("token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

// validateToken validates JWT token
func (wm *WebSocketAuthMiddleware) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(wm.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// extractUserInfo extracts user info from JWT claims
func (wm *WebSocketAuthMiddleware) extractUserInfo(claims jwt.MapClaims) (string, string) {
	userID, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)

	// Fallback to "sub" claim if "user_id" is not available
	if userID == "" {
		userID, _ = claims["sub"].(string)
	}

	return userID, username
}

// RateLimitWebSocket middleware limits number of WebSocket connections
type RateLimitWebSocket struct {
	maxConnections     int
	currentConnections int
	logger             interface {
		Info(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
	}
}

// NewRateLimitWebSocket creates rate limit middleware for WebSocket
func NewRateLimitWebSocket(maxConnections int) *RateLimitWebSocket {
	return &RateLimitWebSocket{
		maxConnections:     maxConnections,
		currentConnections: 0,
		logger:             logging.NewSugaredLogger("websocket-rate-limit"),
	}
}

// RateLimitWebSocketConnection limits number of WebSocket connections
func (rl *RateLimitWebSocket) RateLimitWebSocketConnection(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if rl.currentConnections >= rl.maxConnections {
			rl.logger.Error("WebSocket connection limit exceeded:", rl.currentConnections, "/", rl.maxConnections)
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many connections",
			})
		}

		rl.currentConnections++
		rl.logger.Debug("WebSocket connection established. Current connections:", rl.currentConnections)

		// Decrease counter when connection closes
		c.Response().Header().Set("Connection", "close")
		c.Response().Header().Set("Connection-Close-Callback", "decrease-counter")

		return next(c)
	}
}

// DecreaseConnectionCount decreases connection count
func (rl *RateLimitWebSocket) DecreaseConnectionCount() {
	rl.currentConnections--
	if rl.currentConnections < 0 {
		rl.currentConnections = 0
	}
	rl.logger.Debug("WebSocket connection closed. Current connections:", rl.currentConnections)
}

// CORSWebSocket middleware handles CORS for WebSocket
type CORSWebSocket struct {
	allowedOrigins []string
	logger         interface {
		Info(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
	}
}

// NewCORSWebSocket creates CORS middleware for WebSocket
func NewCORSWebSocket(allowedOrigins []string) *CORSWebSocket {
	return &CORSWebSocket{
		allowedOrigins: allowedOrigins,
		logger:         logging.NewSugaredLogger("websocket-cors"),
	}
}

// CORSWebSocketConnection handles CORS for WebSocket connections
func (cors *CORSWebSocket) CORSWebSocketConnection(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		origin := c.Request().Header.Get("Origin")

		// Check if origin is allowed
		if cors.isOriginAllowed(origin) {
			c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			cors.logger.Error("WebSocket connection from disallowed origin:", origin)
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Origin not allowed",
			})
		}

		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		return next(c)
	}
}

// isOriginAllowed checks if origin is allowed
func (cors *CORSWebSocket) isOriginAllowed(origin string) bool {
	if len(cors.allowedOrigins) == 0 {
		return true // Allow all origins if none specified
	}

	for _, allowedOrigin := range cors.allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}

// WebSocketContextMiddleware middleware to add context for WebSocket
func WebSocketContextMiddleware(logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Add logger to context
			if logger, ok := logger.(*zap.SugaredLogger); ok {
				ctx := logging.AddLoggerToContext(c.Request().Context(), logger)
				c.SetRequest(c.Request().WithContext(ctx))
			}

			// Add request ID for tracking
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = "ws-" + generateRequestID()
			}
			c.Set("request_id", requestID)
			c.Response().Header().Set("X-Request-ID", requestID)

			return next(c)
		}
	}
}

// generateRequestID generates request ID
func generateRequestID() string {
	// Simple implementation - in production, use proper UUID
	return "req-" + string(rune(len("websocket")))
}

// WebSocketLoggingMiddleware middleware to log WebSocket connections
func WebSocketLoggingMiddleware(logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("WebSocket connection attempt from:", c.Request().RemoteAddr, "Path:", c.Request().URL.Path)

			err := next(c)

			if err != nil {
				logger.Error("WebSocket connection failed:", err)
			} else {
				logger.Info("WebSocket connection established successfully")
			}

			return err
		}
	}
}
