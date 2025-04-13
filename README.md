# Quizlet API

A RESTful API for managing quizzes and quiz suites.

## Features

- User authentication and authorization
- CRUD operations for quizzes and quiz suites
- PostgreSQL database with GORM
- Swagger/OpenAPI documentation

## Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your environment variables in `.env`:
   ```
   DB_HOST=localhost
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=quizlet
   DB_PORT=5432
   ```
4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Documentation

The API documentation is available through Swagger UI and ReDoc:

- Swagger UI: `http://localhost:8080/swagger/index.html`
  - Interactive API documentation
  - Test endpoints directly from the browser
  - View request/response schemas
  - See authentication requirements

- OpenAPI Specification: `http://localhost:8080/swagger/doc.json`
  - Raw OpenAPI/Swagger specification in JSON format
  - Useful for generating client libraries or importing into other tools

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   ├── quiz.go
│   │   ├── quiz_suite.go
│   │   └── responses.go
│   └── models/
│       ├── quiz.go
│       └── quiz_suite.go
├── docs/
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Development

### Running Tests

```bash
go test ./...
```

### Database Migrations

Database migrations are handled automatically by GORM when the application starts.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Database Migrations

### Running Migrations in Docker

The project uses `golang-migrate` for database migrations. Migrations can be run directly in the container using `docker-compose exec`:

```bash
# Run all pending migrations
docker-compose exec api migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/quizlet?sslmode=disable" up

# Rollback all migrations
docker-compose exec api migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/quizlet?sslmode=disable" down

# Check migration version
docker-compose exec api migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/quizlet?sslmode=disable" version

# ssh into container
docker-compose exec api /bin/sh
```

Note: Inside the container:
- Database host is `postgres` (service name)
- Database port is `5432` (internal port)
- Migrations are located at `/app/migrations`

### Using Makefile Commands

Alternatively, you can use the Makefile commands:

```bash
# Run migrations up
make migrate-up-container

# Run migrations down
make migrate-down-container
```

## API Endpoints

- `GET /health` - Health check endpoint
- `POST /api/users` - Create a new user
- `GET /api/users/:id` - Get a user by ID
- `PUT /api/users/:id` - Update a user
- `DELETE /api/users/:id` - Delete a user

## Database Connection

PostgreSQL connection details:
- Host: localhost
- Port: 5433
- Username: postgres
- Password: postgres
- Database: quizlet
