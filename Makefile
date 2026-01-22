.PHONY: all build run clean test install deps help

# Build variables
BINARY_NAME=markview
BUILD_DIR=bin
CMD_DIR=cmd/markview
VERSION?=0.1.0

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

all: deps build ## Install dependencies and build the application

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

run-sample: build ## Build and run with sample markdown file
	@echo "Running $(BINARY_NAME) with sample file..."
	./$(BUILD_DIR)/$(BINARY_NAME) -file testdata/sample.md

run-theme: build ## Build and run with theme showcase file
	@echo "Running $(BINARY_NAME) with theme showcase..."
	./$(BUILD_DIR)/$(BINARY_NAME) -file testdata/theme-showcase.md

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps: ## Install/update dependencies
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed"

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run

# Cross-platform builds (requires fyne-cross)
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	fyne-cross linux -arch=amd64,arm64 ./$(CMD_DIR)

build-macos: ## Build for macOS
	@echo "Building for macOS..."
	fyne-cross darwin -arch=amd64,arm64 ./$(CMD_DIR)

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	fyne-cross linux -arch=amd64,arm64 ./$(CMD_DIR)
	fyne-cross darwin -arch=amd64,arm64 ./$(CMD_DIR)

# Fyne packaging
package-macos: build ## Create macOS app bundle
	@echo "Creating macOS app bundle..."
	fyne package -os darwin -icon assets/icon.png

package-linux: build ## Create Linux package
	@echo "Creating Linux package..."
	fyne package -os linux -icon assets/icon.png

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
