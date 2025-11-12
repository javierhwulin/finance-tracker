# Finance Tracker Backend

A clean architecture Go web API using the standard library's `net/http` package and Go 1.22+ ServeMux.

## Structure

```
backend/
├── cmd/
│   └── api/
│       ├── main.go          # Entry point with DI
│       └── main_test.go     # Integration tests
├── internal/
│   ├── app/
│   │   └── app.go           # Composition root
│   ├── domain/
│   │   └── user/
│   │       ├── entity.go    # User entity & validation
│   │       ├── repository.go # Repository interface
│   │       └── service.go   # Domain logic
│   ├── repo/
│   │   └── user_memory.go   # In-memory implementation
│   ├── http/
│   │   ├── router.go        # Routes & middleware
│   │   ├── handler_user.go  # User handlers
│   │   └── dto/
│   │       └── user_dto.go  # Request/Response DTOs
│   └── config/
│       └── config.go        # Configuration
├── go.mod
└── README.md
```

## Quick Start

```bash
# Run server
go run cmd/api/main.go

# Run with custom port
go run cmd/api/main.go -port 3000

# Run tests
go test ./... -v
```

## API Endpoints

### Health Check
```bash
GET /api/health
```

Response:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

### Create User
```bash
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

Response (201 Created):
```json
{
  "id": "uuid-here",
  "email": "user@example.com"
}
```

## Configuration

Set via command-line flags:

- `-port` - Server port (default: 8080)
- `-env` - Environment (default: development)
- `-version` - API version (default: 1.0.0)

## Architecture Highlights

### Clean Architecture
- **Domain Layer**: Pure business logic, no external dependencies
- **Repository Pattern**: Interface in domain, implementation in repo
- **Dependency Injection**: All dependencies wired in `main.go`

### Key Features
- ✅ Standard library only (no external router frameworks)
- ✅ Go 1.22+ ServeMux with method routing
- ✅ Proper dependency injection
- ✅ Thread-safe in-memory storage
- ✅ Request validation (go-playground/validator)
- ✅ Structured logging (slog)
- ✅ Comprehensive tests

## Development Roadmap

### Phase 1: ✅ Complete
- [x] Health check endpoint
- [x] Configuration system
- [x] Structured logging
- [x] Basic tests

### Phase 2: ✅ Complete
- [x] User domain model
- [x] In-memory repository
- [x] Create user endpoint
- [x] Request validation
- [x] Integration tests

### Phase 3: Next Steps
- [ ] PostgreSQL repository implementation
- [ ] Database migrations
- [ ] Authentication (JWT)
- [ ] Password hashing (bcrypt)
- [ ] Middleware (logging, recovery, CORS)
- [ ] More user endpoints (GET, UPDATE, DELETE)

## Notes

- **Passwords**: Currently stored in plaintext. TODO: Implement bcrypt hashing
- **Database**: Using in-memory storage. Ready to swap for Postgres
- **Error Handling**: Basic errors. TODO: Custom error types for better API responses

