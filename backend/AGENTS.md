# AGENTS.md - Development Guidelines for auth_from_scratch

This document provides comprehensive guidelines for agentic coding assistants working on the auth_from_scratch Go project. Follow these conventions to maintain code quality, security, and consistency.

## Build, Test, and Lint Commands

### Building
```bash
# Build the main API binary
go build -o api ./cmd/api

# Build with specific architecture
GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# Build with optimizations
go build -ldflags="-s -w" -o api ./cmd/api
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/app/api/auth

# Run a single test function
go test -run TestFunctionName ./internal/app/api/auth

# Run tests with race detection
go test -race ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting and Formatting
```bash
# Format code (gofmt)
gofmt -w .

# Check formatting without changes
gofmt -d .

# Run go vet (static analysis)
go vet ./...

# Run golint (if installed)
golint ./...

# Run goimports (format imports)
goimports -w .

# Run all checks together
go fmt ./... && go vet ./... && goimports -w .
```

### Database Operations
```bash
# Run database migrations (requires DSN env var)
goose -dir "./internal/app/db/migrations" postgres up

# Create new migration
goose -dir "./internal/app/db/migrations" create migration_name sql

# Rollback last migration
goose -dir "./internal/app/db/migrations" postgres down

# Reset all migrations
goose -dir "./internal/app/db/migrations" postgres reset
```

### Docker Operations
```bash
# Start all services
cd build/ && docker-compose up

# Start in background
cd build/ && docker-compose up -d

# Start only database
cd build/ && docker-compose up -d postgres

# Stop all services
cd build/ && docker-compose down

# View logs
cd build/ && docker-compose logs -f
```

## Code Style Guidelines

### Project Structure
Follow the standard Go project layout:
- `cmd/` - Main applications
- `internal/` - Private application code
- `internal/app/` - Application business logic
- `internal/pkg/` - Shared packages
- `internal/app/db/` - Database layer
- `build/` - Docker and deployment files

### Package Organization
- **Handlers**: HTTP request/response handling
- **Services**: Business logic and orchestration
- **Repos**: Data access layer
- **Types**: Shared type definitions
- **Middleware**: HTTP middleware components

### Import Organization
```go
import (
    // Standard library imports
    "context"
    "crypto/rand"
    "encoding/json"
    "log"
    "time"

    // Third-party imports (alphabetically)
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "github.com/jackc/pgx/v5"

    // Internal imports
    "auth_from_scratch/internal/app/api/auth"
    "auth_from_scratch/internal/pkg/types"
)
```

### Naming Conventions

#### Variables and Functions
- **Exported**: PascalCase (e.g., `UserId`, `ValidateToken`)
- **Unexported**: camelCase (e.g., `userId`, `validateToken`)
- **Constants**: SCREAMING_SNAKE_CASE (e.g., `ACCESS_TOKEN_SECRET`)
- **Acronyms**: Keep consistent (e.g., `userID` not `userId`)

#### Types and Interfaces
```go
// Interfaces: PascalCase with descriptive names
type AuthService interface {
    ValidateToken(ctx context.Context, token string) error
}

// Structs: PascalCase
type AuthHandler struct {
    Service AuthService
}

// Constructors: New + TypeName
func NewAuthHandler(service AuthService) AuthHandler {
    return &AuthHandler{Service: service}
}
```

### Error Handling
```go
// Always handle errors explicitly
result, err := service.ValidateToken(ctx, token)
if err != nil {
    log.Printf("Token validation failed: %v", err)
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "Invalid token",
    })
}

// Use context timeouts
ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
defer cancel()

// Return meaningful error messages to clients
return ah.JSONResponse(c, fiber.StatusBadRequest, fiber.Map{
    "message": "Invalid request format",
    "details": validationErrors,
})
```

### Context and Timeouts
```go
// Use context with appropriate timeouts
func (h *Handler) ProcessRequest(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
    defer cancel()

    // Use ctx for all operations
    result, err := h.service.Process(ctx, data)
    // ...
}
```

### Security Practices

#### Token Handling
- Always hash sensitive tokens before storage
- Use HMAC-SHA256 for token hashing
- Clear cookies on logout/invalidation
- Set secure cookie flags (HttpOnly, Secure, SameSite)

#### Password Security
```go
// Use bcrypt for password hashing
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verify passwords
err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
```

#### Input Validation
- Validate all user inputs
- Use struct tags for JSON validation
- Sanitize database queries with parameterized statements
- Validate OAuth state parameters

### Database Operations
```go
// Use parameterized queries
sql := "SELECT id, email FROM users WHERE id = $1"
row := db.QueryRow(ctx, sql, userId)

// Handle database errors appropriately
if err := row.Scan(&user.Id, &user.Email); err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, errors.New("user not found")
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

### HTTP Response Patterns
```go
// Consistent JSON responses
func (h *Handler) JSONResponse(c *fiber.Ctx, status int, data fiber.Map) error {
    return c.Status(status).JSON(data)
}

// Cookie handling
c.Cookie(&fiber.Cookie{
    Name:     "access_token",
    Value:    token,
    Path:     "/",
    HttpOnly: true,
    Secure:   true,
    SameSite: "Lax",
    Expires:  time.Now().Add(AccessTokenExpiration),
})
```

### Logging
```go
// Use structured logging
log.Printf("User %s authenticated via %s", userId, provider)

// Log security events
log.Printf("Token reuse detected for user %s, session %s", userId, sessionId)

// Include context in error logs
log.Printf("Failed to validate token: %v", err)
```

### Testing Guidelines
```go
// Table-driven tests
func TestValidateToken(t *testing.T) {
    tests := []struct {
        name     string
        token    string
        wantErr  bool
        wantUser string
    }{
        {"valid token", "valid.jwt.token", false, "user123"},
        {"invalid token", "invalid.token", true, ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}

// Mock dependencies for unit tests
type mockAuthService struct {
    validateFunc func(ctx context.Context, token string) error
}

func (m *mockAuthService) ValidateToken(ctx context.Context, token string) error {
    return m.validateFunc(ctx, token)
}
```

### Dependency Injection
```go
// Use interfaces for testability
type AuthHandler interface {
    LoginHandler(c *fiber.Ctx) error
    LogoutHandler(c *fiber.Ctx) error
}

type authHandler struct {
    service AuthService
    repo    AuthRepo
}

// Constructor with dependencies
func NewAuthHandler(service AuthService, repo AuthRepo) AuthHandler {
    return &authHandler{
        service: service,
        repo:    repo,
    }
}
```

### Middleware Usage
```go
// Authentication middleware
func (am *AuthMiddleware) Run() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Cookies("access_token")
        if token == "" {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }
        // validate token...
        return c.Next()
    }
}

// Apply middleware to routes
authRouter.Get("/protected", am.Run(), handler.ProtectedHandler)
```

### Environment Variables
- Use `.env` files for local development
- Prefix with descriptive names (e.g., `GOOGLE_OAUTH_CLIENT_ID`)
- Document required environment variables
- Never commit secrets to repository

### Code Comments
- Document exported functions and types
- Explain complex business logic
- Include security considerations
- Update comments when code changes

Remember: This codebase handles sensitive authentication data. Always prioritize security, input validation, and proper error handling over convenience.</content>
<parameter name="filePath">/home/ynml/code/projects/auth_from_scratch/AGENTS.md