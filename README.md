# Go Chi REST API

Production-ready RESTful API boilerplate built with Go, the lightweight go-chi router, and jackc/pgx/v5 for high-performance PostgreSQL database operations.

---

## 🚀 Features

- **Chi Router:** Idiomatic and performant HTTP routing with grouped API versioning (`/v1`).
- **PostgreSQL + pgx:** High-performance database operations with connection pooling via `pgxpool`.
- **Soft & Hard Deletion:** Supports both soft-deletes (`deleted_at`) and permanent hard-deletes.
- **JSON Serialization:** Pre-configured JSON tags hiding sensitive fields like passwords (`json:"-"`).

---

## 🗄️ Database Schema

### `users` Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

---

## 🛠️ Data Model

```go
type User struct {
    ID        uint       `json:"id"`
    Name      string     `json:"name"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    Password  string     `json:"-"`
    Created_at time.Time  `json:"created_at"`
    Updated_at *time.Time `json:"updated_at"`
    Deleted_at *time.Time `json:"deleted_at"`
}
```

---

## 📋 API Endpoints

All routes are prefixed with /v1

| Methods | Endpoint | Description
| GET | /v1/ping | Health check endpoint
| GET | /v1/users | Retrieve all users
| GET | /v1/users/{id} | Retrieve a user by ID
| POST | /v1/users | Create a new user
| PUT | /v1/users/{id} | Update an existing user
| DELETE | /v1/users/{id} | Soft delete a user by ID
| DELETE | /v1/users-delete/{id} | Permanently delete a user by ID

---

## ⚙️ Getting Started

### Prerequisites

- Go installed (version 1.22 or higher)
- PostgreSQL instance running locally or remotely

### Installation

1. Clone the repository:

```bash
git clone https://github.com/AlifRJ/Go_Auth.git
cd Go_Auth
```

2. Install dependencies:

```bash
go mod download
```

3. Configure environment variables

```bash
DB_DATABASE={postgres}
DB_PORT={5432}
DB_USER={postgres}
DB_PASSWORD={password}
```

4. Run the application:

```bash
go run main.go
```

---

## 📮 API Payload Examples

### Health Check (GET /v1/ping)

```bash
PONG!
```

### Create User (POST /v1/users)

Request Body:

```json
{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "password": "secretpassword"
}
```

Response (201 Created):

```json
{
  "id": 1,
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2026-09-07T16:30:00Z",
  "updated_at": null,
  "deleted_at": null
}
```

### Update User (PUT /v1/users/1)

Request Body:

```json
{
  "name": "Jane H. Doe",
  "username": "janedoe_updated",
  "email": "jane.new@example.com"
}
```

Response (201 Created):

```json
{
  "id": 1,
  "name": "Jane H. Doe",
  "username": "janedoe_updated",
  "email": "jane.new@example.com",
  "created_at": "2026-09-07T16:30:00Z",
  "updated_at": "2026-09-07T16:35:00Z",
  "deleted_at": null
}
```

<FollowUp label="Want me to generate full CRUD handler implementations for Chi and pgx?" query="Generate the full Go HTTP handler implementations for CRUD operations on this User model using Chi router and pgx."/>
