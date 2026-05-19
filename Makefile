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
	@echo "Running tests..."
	go test ./... -v -race -count=1

lint:
	@echo "Running linter..."
	golangci-lint run ./...

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
