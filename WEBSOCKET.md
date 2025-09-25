# WebSocket Integration

This document describes the WebSocket functionality integrated into the user management API.

## Features

- Real-time bidirectional communication
- Room-based messaging
- User authentication support
- Connection management
- Message broadcasting
- Heartbeat/ping-pong mechanism

## WebSocket Endpoints

### 1. Basic WebSocket Connection
```
GET /ws
```
Connects to WebSocket without authentication.

### 2. WebSocket with Room
```
GET /ws/{room}
```
Connects to WebSocket and joins a specific room.

### 3. Authenticated WebSocket Connection
```
GET /ws/auth
```
Connects to WebSocket with JWT authentication.

### 4. Authenticated WebSocket with Room
```
GET /ws/auth/{room}
```
Connects to WebSocket with authentication and joins a room.

## Authentication

For authenticated endpoints, include the JWT token in one of these ways:

1. **Authorization Header:**
   ```
   Authorization: Bearer <jwt_token>
   ```

2. **Query Parameter:**
   ```
   /ws/auth?token=<jwt_token>
   ```

3. **Cookie:**
   ```
   Cookie: token=<jwt_token>
   ```

## Message Format

All WebSocket messages follow this JSON structure:

```json
{
  "type": "message_type",
  "room": "room_name",
  "content": "message_content",
  "user_id": "user_id",
  "username": "username",
  "timestamp": 1234567890,
  "data": {}
}
```

## Message Types

### 1. Join Room
```json
{
  "type": "join",
  "room": "room_name"
}
```

### 2. Leave Room
```json
{
  "type": "leave",
  "room": "room_name"
}
```

### 3. Send Message
```json
{
  "type": "message",
  "room": "room_name",
  "content": "Hello, world!"
}
```

### 4. Ping
```json
{
  "type": "ping"
}
```

### 5. Pong
```json
{
  "type": "pong"
}
```

### 6. User List
```json
{
  "type": "user_list",
  "room": "room_name",
  "data": ["user1", "user2", "user3"]
}
```

### 7. Error
```json
{
  "type": "error",
  "content": "Error message"
}
```

## Client Example (JavaScript)

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws/room1');

// Send authentication token (if required)
ws.onopen = function() {
    // Send join message
    ws.send(JSON.stringify({
        type: 'join',
        room: 'room1'
    }));
};

// Handle incoming messages
ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    
    switch(message.type) {
        case 'message':
            console.log(`${message.username}: ${message.content}`);
            break;
        case 'user_list':
            console.log('Users in room:', message.data);
            break;
        case 'join':
            console.log(`${message.username} joined the room`);
            break;
        case 'leave':
            console.log(`${message.username} left the room`);
            break;
    }
};

// Send a message
function sendMessage(content) {
    ws.send(JSON.stringify({
        type: 'message',
        room: 'room1',
        content: content
    }));
}

// Handle connection close
ws.onclose = function() {
    console.log('WebSocket connection closed');
};

// Handle errors
ws.onerror = function(error) {
    console.error('WebSocket error:', error);
};
```

## Configuration

WebSocket behavior can be configured using environment variables:

```bash
# Enable/disable WebSocket
WEBSOCKET_ENABLED=true

# Maximum concurrent connections
WEBSOCKET_MAX_CONNECTIONS=1000

# Buffer sizes
WEBSOCKET_READ_BUFFER_SIZE=1024
WEBSOCKET_WRITE_BUFFER_SIZE=1024

# Message limits
WEBSOCKET_MAX_MESSAGE_SIZE=512

# Timeouts
WEBSOCKET_WRITE_WAIT=10s
WEBSOCKET_PONG_WAIT=60s
WEBSOCKET_PING_PERIOD=54s

# CORS
WEBSOCKET_ALLOWED_ORIGINS=*

# Authentication
WEBSOCKET_REQUIRE_AUTH=false

# Heartbeat
WEBSOCKET_HEARTBEAT_INTERVAL=30s
WEBSOCKET_CONNECTION_TIMEOUT=60s
```

## Room Management

- Users can join multiple rooms
- Messages are broadcast to all users in the same room
- Room membership is managed automatically
- Users are notified when someone joins/leaves a room

## Connection Management

- Automatic cleanup of disconnected clients
- Heartbeat mechanism to detect dead connections
- Rate limiting for connection attempts
- CORS support for web clients

## Security Considerations

- Use HTTPS/WSS in production
- Implement proper JWT token validation
- Set appropriate CORS origins
- Monitor connection limits
- Implement rate limiting for messages

## Error Handling

Common error scenarios:

1. **Authentication Failed**: Invalid or missing JWT token
2. **Room Full**: Maximum connections reached
3. **Invalid Message**: Malformed JSON or unsupported message type
4. **Connection Timeout**: Client not responding to ping

## Monitoring

Monitor WebSocket connections using:

- Connection count metrics
- Message throughput
- Error rates
- Room occupancy
- Client connection duration
