# Go Authentication API

A Go REST API for user registration, JWT authentication, refresh-token rotation, logout, and user management. The service uses PostgreSQL for durable data, Redis for caching and token revocation, and Chi for HTTP routing.

## Features

- Go HTTP API with Chi router and CORS support
- User registration and login with either username or email
- Short-lived access JWTs and seven-day refresh JWTs
- HTTP-only refresh-token cookie with refresh-token rotation
- Logout that removes the persisted refresh token and blacklists the access token in Redis
- PostgreSQL persistence with automatic startup migrations
- Redis-backed user caching and token blacklist checks
- Password hashing with `bcrypt`
- Soft delete and permanent delete support
- Clean separation across handler, service, repository, model, middleware, and migration layers

---

## Tech Stack

- Go
- Chi Router
- PostgreSQL (`pgxpool`)
- Redis (`go-redis`)
- pgx
- JWT (`github.com/go-chi/jwtauth/v5`)
- bcrypt
- godotenv

---

## Project Structure

```text
.
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── README.md
├── app/
│   ├── handler/
│   │   ├── authHandler.go
│   │   └── userHandler.go
│   ├── middleware/
│   │   └── blacklist.go
│   ├── migration/
│   │   └── userMigration.go
│   ├── model/
│   │   ├── auth.go
│   │   └── user.go
│   ├── repository/
│   │   ├── authRepository.go
│   │   ├── cachedAuthRepository.go
│   │   ├── cachedUserRepository.go
│   │   └── userRepository.go
│   └── service/
│       ├── authService.go
│       └── userService.go
└──
```

---

## Database Schema

The application creates the `users` and `auth_sessions` tables, plus lookup indexes, automatically at startup if they do not already exist.

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth_sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

---

## Environment Variables

Create a `.env` file, or provide these values through the environment:

```env
APP_PORT=8080
DB_HOST=localhost
DB_DATABASE=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
ACCESS_SECRET_KEY=your-access-token-secret
REFRESH_SECRET_KEY=your-refresh-token-secret
```

`ACCESS_SECRET_KEY` and `REFRESH_SECRET_KEY` are required. The application exits during startup if PostgreSQL, Redis, or either JWT secret is unavailable.

---

## Quick Start

### Prerequisites

- Go 1.26.5+
- PostgreSQL installed and running
- Redis installed and running
- A PostgreSQL database and Redis instance accessible using the configured host and port

### Installation

```bash
go mod download
```

### Run the API

```bash
go run main.go
```

The server listens on `APP_PORT` (normally `8080`).

---

## API Routes

The health endpoint is available at `/health`. Authentication and user routes are prefixed with `/v1`.

### Public Endpoints

| Method | Endpoint       | Description                                    |
| ------ | -------------- | ---------------------------------------------- |
| POST   | `/v1/register` | Register a user                                |
| POST   | `/v1/login`    | Authenticate and return an access JWT          |
| POST   | `/v1/refresh`  | Rotate tokens using the refresh JWT and cookie |
| GET    | `/health`      | Health check                                   |

### Protected Endpoints

Protected routes require a Bearer token in the `Authorization` header:

```http
Authorization: Bearer <token>
```

| Method | Endpoint                   | Description                     |
| ------ | -------------------------- | ------------------------------- |
| GET    | `/v1/me`                   | Return the authenticated claims |
| POST   | `/v1/logout`               | Revoke the current session      |
| GET    | `/v1/users/`               | List users                      |
| GET    | `/v1/users/{id}`           | Get a user by ID                |
| POST   | `/v1/users/`               | Create a user                   |
| PUT    | `/v1/users/{id}`           | Update a user                   |
| DELETE | `/v1/users/{id}`           | Soft delete a user              |
| DELETE | `/v1/users/{id}/permanent` | Permanently delete a user       |

`GET /v1/users/` accepts optional `limit` and `offset` query parameters.

---

## Example: Login

### Request

```http
POST /v1/login
Content-Type: application/json
```

```json
{
  "identity": "johndoe",
  "password": "secretpassword"
}
```

### Response

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

The response also sets an HTTP-only `refresh_token` cookie scoped to `/v1`. Use the `access_token` as the `Authorization: Bearer <token>` header for protected endpoints.

### Refresh and Logout

Refresh tokens are stored in PostgreSQL and Redis. The refresh endpoint validates the refresh JWT and cookie, then rotates both tokens:

```http
POST /v1/refresh
Authorization: Bearer <refresh-token>
Cookie: refresh_token=<refresh-token>
```

Logout deletes the stored refresh token, blacklists the access-token JTI until it expires, and clears the refresh cookie.

---

## Example: Create User

### Request

```http
POST /v1/users/
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "secretpassword"
}
```

### Response

```json
{
  "id": 1,
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2026-09-10T12:00:00Z",
  "updated_at": null,
  "deleted_at": null
}
```

---

## Authentication and Caching

- Login accepts either a `username` or `email` in the `identity` field.
- Access tokens expire after 15 minutes; refresh tokens expire after seven days.
- Passwords are hashed before being stored in PostgreSQL and are never returned in JSON responses.
- User records loaded by ID are cached in Redis for 15 minutes and invalidated after updates or deletes.
- Protected user routes reject access tokens whose JTI is present in the Redis blacklist.
- Soft delete marks a record with `deleted_at` while preserving the row; permanent delete removes it.

---

## Use Cases

This repository is a practical foundation for authentication-first services, admin APIs, and applications that need PostgreSQL persistence with Redis-backed session and cache support.
