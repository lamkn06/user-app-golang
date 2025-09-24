# User App Golang

A RESTful API for user management built with Go, Echo framework, and PostgreSQL.

## Features

- User CRUD operations
- Input validation
- Error handling with detailed responses
- Logging with structured logging
- Swagger API documentation
- Database migrations
- Docker support

## Tech Stack

- **Framework**: Echo v4
- **Database**: PostgreSQL with PostGIS
- **ORM**: Bun
- **Validation**: go-playground/validator
- **Logging**: Zap
- **Documentation**: Swagger/OpenAPI
- **Migration**: golang-migrate

## Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (optional)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd user-app-golang
go mod download
```

### 2. Environment Setup

Create `.env` file:

```env
PORT=8080
ENVIRONMENT=development
API_VERSION=v1
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=userapp
```

### 3. Database Setup

#### Option A: Using Docker Compose

```bash
docker-compose up -d
```

#### Option B: Manual PostgreSQL Setup

```bash
# Create database
createdb userapp

# Run migrations
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/userapp?sslmode=disable" up
```

### 4. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Swagger UI

Access the interactive API documentation at:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

### API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID

## Swagger Documentation

### Generating Swagger Docs

1. **Install Swagger CLI**:
```bash
go install github.com/swaggo/swag/cmd/swag@v1.8.12
```

2. **Generate Documentation**:
```bash
swag init -g main.go
```

This will create:
- `docs/docs.go` - Generated Go code
- `docs/swagger.json` - JSON specification
- `docs/swagger.yaml` - YAML specification

3. **Update Documentation**:
Whenever you add new endpoints or modify existing ones, regenerate the docs:
```bash
swag init -g main.go
```

### Adding Swagger Annotations

Add comments to your handlers:

```go
// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.NewUserRequest true "User information"
// @Success 200 {object} response.NewUserResponse
// @Failure 400 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /users [post]
func (r *UserRouter) CreateUser(c echo.Context) error {
    // handler implementation
}
```

## Project Structure

```
├── cmd/                    # Application entrypoints
├── internal/               # Private application code
│   ├── middleware/         # HTTP middleware
│   ├── repository/         # Data access layer
│   ├── route/             # HTTP handlers
│   ├── runtime/           # Configuration
│   └── service/           # Business logic
├── migrations/            # Database migrations
├── pkg/                   # Public packages
│   ├── api/              # Request/Response models
│   ├── exception/        # Error handling
│   └── logging/          # Logging utilities
├── docs/                 # Generated Swagger docs
├── docker-compose.yml    # Docker services
└── main.go              # Application entry point
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## Docker

### Build Image

```bash
docker build -t user-app-golang .
```

### Run with Docker Compose

```bash
docker-compose up -d
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `API_VERSION` | API version | `v1` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `userapp` |

## Error Handling

The API returns structured error responses:

```json
{
  "code": "VALIDATION",
  "message": "Validation failed",
  "details": [
    {
      "key": "NewUserRequest.Email",
      "field": "Email",
      "message": "Failed on the 'required' tag"
    }
  ]
}
```

## Logging

The application uses structured logging with Zap:

- **Development**: Pretty console output with colors
- **Production**: JSON format for log aggregation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update Swagger documentation
6. Submit a pull request

## License

This project is licensed under the MIT License.
