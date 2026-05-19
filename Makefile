MODULE := github.com/rizky/go-scaffold
BUILD_DIR := bin

build:
	@echo "Building all services..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/api ./cmd/api
	go build -o $(BUILD_DIR)/worker ./cmd/worker
	@echo "Build complete: $(BUILD_DIR)/"

run-api:
	@echo "Starting API server..."
	APP_NAME=$(MODULE) APP_ENV=development go run ./cmd/api

run-worker:
	@echo "Starting Worker..."
	APP_NAME=$(MODULE) APP_ENV=development go run ./cmd/worker

dev:
	@echo "Starting all services..."
	docker compose -f deploy/docker-compose.yml up --build

dev-down:
	docker compose -f deploy/docker-compose.yml down

dev-logs:
	docker compose -f deploy/docker-compose.yml logs -f

test:
	@echo "Running all tests (short mode — skips integration)..."
	go test -short -v -race -count=1 ./...

test-integration:
	@echo "Running integration tests (requires Docker)..."
	go test -run Integration -v -race -count=1 ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -short -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run ./...

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
