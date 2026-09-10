# Go REST API Boilerplate

A ready-to-use Go REST API starter for building secure, production-style backend services with PostgreSQL, JWT authentication, and a clean layered architecture.

This project is designed to be a solid foundation for a full SaaS backend, admin panel API, or internal service.

## Features

- Go HTTP API with Chi router
- JWT-based authentication and protected routes
- PostgreSQL database integration using `pgxpool`
- User management with CRUD-style endpoints
- Password hashing with `bcrypt`
- Soft delete and permanent delete support
- Environment-based configuration via `.env`
- Clean separation of concerns across handler, service, repository, model, and migration layers

---

## Tech Stack

- Go
- Chi Router
- PostgreSQL
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
├── .env-example
├── README.md
├── app/
│   ├── handler/
│   │   └── userHandler.go
│   ├── migration/
│   │   └── userMigration.go
│   ├── model/
│   │   └── user.go
│   ├── repository/
│   │   └── userRepo.go
│   └── service/
│       └── userService.go
└──
```

---

## Database Schema

The application creates the `users` table automatically if it does not already exist.

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
```

---

## Environment Variables

Create a `.env` file using the example below:

```env
DB_DATABASE=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
SECRET_KEY=your-random-64-character-secret-key
```

A sample is available in `.env-example`.

---

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL installed and running
- A local database accessible on `localhost`

### Installation

```bash
go mod download
```

### Run the API

```bash
go run main.go
```

The server starts on port `8080`.

---

## API Routes

All routes are prefixed with `/v1`.

### Public Endpoints

| Method | Endpoint    | Description                          |
| ------ | ----------- | ------------------------------------ |
| GET    | `/v1/ping`  | Health check                         |
| POST   | `/v1/login` | Authenticate a user and return a JWT |

### Protected Endpoints

Protected routes require a Bearer token in the `Authorization` header:

```http
Authorization: Bearer <token>
```

| Method | Endpoint                | Description               |
| ------ | ----------------------- | ------------------------- |
| GET    | `/v1/users`             | List all users            |
| GET    | `/v1/users/{id}`        | Get user by ID            |
| POST   | `/v1/users`             | Create a new user         |
| PUT    | `/v1/users/{id}`        | Update a user             |
| DELETE | `/v1/users/{id}`        | Soft delete a user        |
| DELETE | `/v1/users-delete/{id}` | Permanently delete a user |

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Use the returned token as the `Authorization: Bearer <token>` header for protected endpoints.

---

## Example: Create User

### Request

```http
POST /v1/users
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

## Notes

- Login accepts either a `username` or `email` in the `identity` field.
- Passwords are hashed before being stored in the database.
- Soft delete marks the record with `deleted_at` while preserving the row.
- Permanent delete removes the record from the database.

---

## Why Use This Boilerplate?

This repository is a practical starting point for developers who want a clean, secure, and extensible Go API foundation without spending time setting up the basic architecture from scratch.

It is especially useful for:

- startup MVPs
- internal services
- admin APIs
- authentication-first applications
- learning and experimentation
