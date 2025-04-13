# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies and migrate tool
RUN apk add --no-cache gcc musl-dev curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and migrations
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/api/main.go && \
    ls -la docs/

# Build the application
RUN cd cmd/api && CGO_ENABLED=0 GOOS=linux go build -o ../../main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install migrate tool in final stage
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Copy the binary, migrations, and docs from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Verify files are copied correctly
RUN ls -la docs/

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"] 