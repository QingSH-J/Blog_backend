# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a comprehensive Go-based backend API built with Gin framework. The application provides multiple features including user authentication, book logging/reviews, forum discussions, reading time tracking, and AI-powered chat functionality.

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
- `DB_DSN` - PostgreSQL connection string (required)
- `JWT_SECRET` - JWT signing secret (required)
- `OPENAI_API_KEY` - OpenAI/DeepSeek API key for chat functionality (required)
- `SERVER_HOST` - Server host (default: localhost)
- `SERVER_PORT` - Server port (default: 8080)
- `JWT_EXPIRES_IN` - JWT token expiration (default: 24h)
- `ENVIRONMENT` - Environment mode (default: development)
- `LOG_LEVEL` - Logging level (default: info)
- `CORS_ALLOWED_ORIGINS` - Allowed CORS origins (default: http://localhost:3000)
- `CORS_ALLOWED_METHODS` - Allowed HTTP methods (default: GET,POST,PUT,DELETE,OPTIONS)
- `CORS_ALLOWED_HEADERS` - Allowed headers (default: Content-Type,Authorization)
- `EXTERNAL_API_BASE_URL` - External API base URL
- `EXTERNAL_API_KEY` - External API key
- `BCRYPT_COST` - Password hashing cost (default: 12)

## Architecture

### Project Structure
- `main.go` - Application entry point with dependency injection
- `internal/config/` - Environment configuration management
- `internal/model/` - GORM data models (UserLog, BookLog, Topic, Comment, ReadTime, ChatLog, Message)
- `internal/store/` - Database layer (repositories)
- `internal/service/` - Business logic layer
- `internal/api/` - HTTP handlers and routing
- `internal/api/middleware/` - HTTP middleware (auth, CORS)
- `internal/external/` - External API integrations

### Key Architecture Patterns

**Dependency Injection**: Services are injected into handlers via `HandlerDependencies` struct in `api/router.go:11-17`

**Database Layer**: 
- Singleton database connection in `store/database.go`
- Store pattern with dedicated stores for all entities
- Auto-migrations run on application startup for all models

**Authentication**:
- JWT-based authentication with Bearer tokens
- Auth middleware validates tokens and sets user context at `middleware/auth_middleware.go`
- Protected routes use `middleware.AuthMiddleware()`

**API Routes**:
- `/api/v1/auth/` - User registration and login (public)
- `/api/v1/review/books` - Get user's book library (protected)
- `/api/v1/new/` - Create book entries (protected)
- `/api/v1/books/:id` - Get/update specific books (protected)
- `/api/v1/search` - Book search functionality (public)
- `/api/v1/forum/` - Forum topics and comments (mixed public/protected)
- `/api/v1/readtime/` - Reading time tracking (protected)
- `/api/v1/chat/` - AI chat functionality (protected)

### Database Models
- `UserLog`: Users with email/password authentication
- `BookLog`: Book entries with external metadata + user ratings/comments
- `Topic`: Forum discussion topics
- `Comment`: Comments on forum topics
- `ReadTime`: User reading time tracking records
- `ChatLog`: AI chat sessions
- `Message`: Individual messages within chat sessions
- Foreign key relationships with CASCADE constraints for data integrity

### Key Features

**Book Management**:
- Personal book library with ratings and comments
- Book search functionality
- Integration with external book APIs

**Forum System**:
- Topic creation and discussion
- Comment system with user attribution
- View count tracking

**Reading Time Tracking**:
- Record reading sessions
- Weekly reading time analytics

**AI Chat Integration**:
- OpenAI/DeepSeek API integration for intelligent conversations
- Chat session management
- Message history tracking

### Configuration
Centralized config management in `internal/config/config.go` with singleton pattern and comprehensive environment variable support with sensible defaults.