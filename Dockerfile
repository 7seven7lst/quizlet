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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install migrate tool in final stage
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Copy the binary and migrations from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"] 