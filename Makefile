.PHONY: build test clean run-example docker-up docker-down

# Build the library
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean ./...
	rm -rf bin/ dist/

# Run basic SSH example
run-example:
	go run examples/basic_ssh/main.go

# Start guacd with Docker Compose
docker-up:
	docker-compose up -d

# Stop guacd
docker-down:
	docker-compose down

# Run all: start docker and run example
run-all: docker-up
	@echo "Waiting for guacd to start..."
	@sleep 5
	@echo "Running example..."
	@make run-example
