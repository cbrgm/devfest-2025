.PHONY: build clean run test help install-config

# Binary name
BINARY_NAME=teashop-mcp-server

# Build directory
BUILD_DIR=bin

# Claude Desktop config path
CLAUDE_CONFIG_DIR=~/Library/Application Support/Claude
CLAUDE_CONFIG_FILE=$(CLAUDE_CONFIG_DIR)/claude_desktop_config.json

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(GOBIN)/$(BINARY_NAME) .
	@echo "Binary built: $(GOBIN)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean completed"

# Run the server
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(GOBIN)/$(BINARY_NAME)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	@go install .
	@echo "Installation completed"

# Install MCP config to Claude Desktop
install-config: build
	@echo "Installing MCP config to Claude Desktop..."
	@mkdir -p "$(CLAUDE_CONFIG_DIR)"
	@if [ -f "$(CLAUDE_CONFIG_FILE)" ]; then \
		echo "Backing up existing config to $(CLAUDE_CONFIG_FILE).backup"; \
		cp "$(CLAUDE_CONFIG_FILE)" "$(CLAUDE_CONFIG_FILE).backup"; \
	fi
	cp mcp-config.json "$(CLAUDE_CONFIG_FILE)"
	@echo "MCP config installed to $(CLAUDE_CONFIG_FILE)"
	@echo "Restart Claude Desktop to use the tea shop MCP server"

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run the server"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  test          - Run tests"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  install-config - Install MCP config to Claude Desktop"
	@echo "  help          - Show this help message"