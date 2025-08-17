# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based book logging/review backend API built with Gin framework. The application allows users to register, login, and maintain a personal book library with ratings and comments.

## Development Commands

### Build and Run
```bash
go mod tidy                # Install/update dependencies
go run main.go            # Run the application
go build -o app main.go   # Build binary
```

### Database
The application uses PostgreSQL with GORM for migrations and queries. Database migrations run automatically on startup via the `Migrate()` methods in stores.

### Environment Setup
Create a `.env` file with required environment variables:
- `DB_DSN` - PostgreSQL connection string
- `JWT_SECRET` - JWT signing secret (required)
- `SERVER_PORT` - Server port (default: 8080)
- `CORS_ALLOWED_ORIGINS` - Allowed CORS origins (default: http://localhost:3000)

## Architecture

### Project Structure
- `main.go` - Application entry point with dependency injection
- `internal/config/` - Environment configuration management
- `internal/model/` - GORM data models (UserLog, BookLog)
- `internal/store/` - Database layer (repositories)
- `internal/service/` - Business logic layer
- `internal/api/` - HTTP handlers and routing
- `internal/api/middleware/` - HTTP middleware (auth, CORS)

### Key Architecture Patterns

**Dependency Injection**: Services are injected into handlers via `HandlerDependencies` struct in `api/router.go:11-14`

**Database Layer**: 
- Singleton database connection in `store/database.go:25-37`
- Store pattern with dedicated stores for users and book logs
- Auto-migrations run on application startup

**Authentication**:
- JWT-based authentication with Bearer tokens
- Auth middleware validates tokens and sets user context at `middleware/auth_middleware.go:13`
- Protected routes use `middleware.AuthMiddleware()`

**API Routes**:
- `/api/v1/auth/` - Registration and login (public)
- `/api/v1/review/books` - Get user's book library (protected)
- `/api/v1/new/` - Create book entries (protected)
- `/api/v1/books/:id` - Get/update specific books (protected)
- `/api/v1/search` - Search functionality (public)

### Database Models
- `UserLog`: Users with email/password authentication
- `BookLog`: Book entries with external metadata + user ratings/comments
- Foreign key relationship: BookLog.UserID → UserLog.ID with CASCADE

### Configuration
Centralized config management in `internal/config/config.go` with singleton pattern and environment variable defaults.