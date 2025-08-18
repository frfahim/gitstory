# GitStory Makefile

.PHONY: build clean test install run help lint fmt release tag check-version security-scan

# Build variables
BINARY_NAME=gitstory
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Default target
all: test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-arm64.exe .
	@echo "Built binaries:"
	@ls -la bin/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install the application
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .

# Run the application (development)
run:
	go run . $(ARGS)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache
	go clean -testcache

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Security scanning
security-scan:
	@echo "Running security scan..."
	@command -v govulncheck >/dev/null 2>&1 || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Validate version format (semantic versioning)
check-version:
	@echo "Checking version format..."
	@echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+' || (echo "Version should follow semantic versioning (e.g., v1.0.0)" && exit 1)

# Create a new release tag
tag:
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version (e.g., v0.1.0): " NEW_VERSION; \
	if [ -z "$$NEW_VERSION" ]; then \
		echo "Version cannot be empty"; \
		exit 1; \
	fi; \
	echo "$$NEW_VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+' || (echo "Version should follow semantic versioning (e.g., v1.0.0)" && exit 1); \
	git tag -a $$NEW_VERSION -m "Release $$NEW_VERSION"; \
	echo "Created tag $$NEW_VERSION. Push with: git push origin $$NEW_VERSION"

# Release workflow
release: test lint
	@echo "Creating release..."
	@read -p "Enter version (e.g., v0.1.0): " NEW_VERSION; \
	if [ -z "$$NEW_VERSION" ]; then \
		echo "Version cannot be empty"; \
		exit 1; \
	fi; \
	echo "$$NEW_VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+' || (echo "Version should follow semantic versioning (e.g., v1.0.0)" && exit 1); \
	echo "Building release binaries..."; \
	$(MAKE) build-all; \
	echo "Creating and pushing tag..."; \
	git tag -a $$NEW_VERSION -m "Release $$NEW_VERSION"; \
	git push origin $$NEW_VERSION; \
	echo "Release $$NEW_VERSION created and pushed!"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download

# Help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Building:"
	@echo "  build          Build the application"
	@echo "  build-all      Build for multiple platforms"
	@echo "  install        Install the application"
	@echo "  clean          Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt            Format code"
	@echo "  lint           Lint code (requires golangci-lint)"
	@echo "  security-scan  Run security vulnerability scan"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps           Download dependencies"
	@echo "  update-deps    Update dependencies"
	@echo "  dev-setup     Setup development environment"
	@echo ""
	@echo "Release Management:"
	@echo "  tag            Create a new version tag"
	@echo "  release        Full release workflow (test, lint, build, tag)"
	@echo "  check-version  Validate version format"
	@echo ""
	@echo "Development:"
	@echo "  run            Run the application (use ARGS=... for arguments)"
	@echo "  help           Show this help message"
