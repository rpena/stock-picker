APP_NAME := stock-picker
BUILD_DIR := build
BINARY_NAME := $(BUILD_DIR)/$(APP_NAME)
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
# Docker image name
IMAGE_NAME := $(APP_NAME)
# Docker image tag
IMAGE_TAG := 1.0
# Dockerfile name
DOCKERFILE := Dockerfile

# Default target: build the application
all: build

# Build the Go application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BINARY_NAME) cmd/main.go
	@echo "Build complete: $(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the Go application
run: build
	@echo "Running $(APP_NAME)..."
	@./$(BINARY_NAME)

# Clean up build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "Formatting complete."

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "Dependencies installed."

docker-build:
	@echo "Building Docker image $(IMAGE_NAME):$(IMAGE_TAG) using $(DOCKERFILE)..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f $(DOCKERFILE) .
	@echo "Docker image $(IMAGE_NAME):$(IMAGE_TAG) built successfully."

docker-run:
	@echo "Running Docker container $(IMAGE_NAME):$(IMAGE_TAG)..."
	docker run -d -p 8080:8080 --name $(APP_NAME)-container $(IMAGE_NAME):$(IMAGE_TAG)
	@echo "Container $(APP_NAME)-container started. Access at http://localhost:8080/stock/"

.PHONY: all build test run clean fmt lint deps docker-build