.PHONY: run build test swagger clean

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/docs.go
	rm -rf docs/swagger.json
	rm -rf docs/swagger.yaml

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run with hot reload (requires air)
dev:
	air

# Database migration
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down



