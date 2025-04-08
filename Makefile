.PHONY: migrate-up migrate-down migrate-create migrate-up-container migrate-down-container

# Database connection details
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=quizlet

# Container database connection details
CONTAINER_DB_HOST=postgres
CONTAINER_DB_PORT=5432
CONTAINER_DB_USER=postgres
CONTAINER_DB_PASSWORD=postgres
CONTAINER_DB_NAME=quizlet

# Local migration commands
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Container migration commands
migrate-up-container:
	docker-compose exec api migrate -path /app/migrations -database "postgres://$(CONTAINER_DB_USER):$(CONTAINER_DB_PASSWORD)@$(CONTAINER_DB_HOST):$(CONTAINER_DB_PORT)/$(CONTAINER_DB_NAME)?sslmode=disable" up

migrate-down-container:
	docker-compose exec api migrate -path /app/migrations -database "postgres://$(CONTAINER_DB_USER):$(CONTAINER_DB_PASSWORD)@$(CONTAINER_DB_HOST):$(CONTAINER_DB_PORT)/$(CONTAINER_DB_NAME)?sslmode=disable" down 