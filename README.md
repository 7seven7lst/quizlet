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

You can run tests in two ways:

1. **Locally**:
   ```bash
   # Run tests
   go test ./...

   # Run tests with coverage report
   go test -coverprofile=coverage.txt -covermode=atomic ./...
   
   # View coverage in browser
   go tool cover -html=coverage.txt
   
   # View coverage in terminal
   go tool cover -func=coverage.txt
   ```

2. **In Docker Container**:
   ```bash
   # Start the development container
   docker-compose up -d api-dev

   # Run all tests
   docker-compose exec api-dev go test ./...

   # Run tests with coverage
   docker-compose exec api-dev go test -cover ./...

   # Run tests with detailed coverage report
   docker-compose exec api-dev go test -coverprofile=coverage.txt -covermode=atomic ./...
   docker-compose exec api-dev go tool cover -func=coverage.txt
   ```

   Note: The development container (`api-dev`) includes the Go toolchain and mounts your local code directory, so any changes to your code or tests will be immediately reflected.

### Continuous Integration

The project uses GitHub Actions for continuous integration. The test suite runs automatically on:
- Every push to main/master branch
- Every pull request to main/master branch

The CI pipeline:
1. Sets up a PostgreSQL database
2. Installs Go and project dependencies
3. Runs the entire test suite with race condition detection
4. Generates and uploads code coverage reports to Codecov

#### Code Coverage Requirements

The project maintains strict code coverage requirements:
- Overall project coverage must be at least 80%
- New code in pull requests must have at least 80% coverage
- Coverage is only measured for code in the `internal` directory
- Test files, mocks, and generated code are excluded from coverage calculations
- Coverage reports are available on [Codecov](https://codecov.io)

You can view:
- Test results in the "Actions" tab of the GitHub repository
- Code coverage reports and trends on Codecov
- Line-by-line coverage annotations in pull requests

To skip CI for a particular commit, include `[skip ci]` in your commit message.

[![Test Suite](https://github.com/{username}/{repo}/actions/workflows/test.yml/badge.svg)](https://github.com/{username}/{repo}/actions)
[![codecov](https://codecov.io/gh/{username}/{repo}/branch/main/graph/badge.svg)](https://codecov.io/gh/{username}/{repo})

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
