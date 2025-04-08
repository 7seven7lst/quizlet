# Quizlet API

A Go-based API service with PostgreSQL database.

## Prerequisites

- Docker
- Docker Compose
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository
2. Build and start the containers:
```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080` and PostgreSQL at `localhost:5433`.

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
