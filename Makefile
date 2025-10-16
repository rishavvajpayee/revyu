.PHONY: build install clean test verify help

# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Variables
BINARY_NAME=revyu
BUILD_DIR=.
INSTALL_DIR=/usr/local/bin
GO=go
LDFLAGS=-ldflags "-X main.buildTimeAPIKey=$(OPENAI_API_KEY)"

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  \033[36mbuild\033[0m           Build the binary with code signing"
	@echo "  \033[36minstall\033[0m         Build and install the binary to /usr/local/bin"
	@echo "  \033[36mverify\033[0m          Verify the code signature of the binary"
	@echo "  \033[36mverify-installed\033[0m Verify the code signature of the installed binary"
	@echo "  \033[36mclean\033[0m           Remove built binaries"
	@echo "  \033[36muninstall\033[0m       Remove the installed binary from /usr/local/bin"
	@echo "  \033[36mrebuild\033[0m         Clean and rebuild"
	@echo "  \033[36mreinstall\033[0m       Clean, rebuild, and reinstall"
	@echo "  \033[36mtest\033[0m            Run tests"
	@echo "  \033[36mhelp\033[0m            Show this help message"

build: ## Build the binary with code signing
	@echo "Building $(BINARY_NAME) with embedded OpenAI API key..."
	@if [ -z "$(OPENAI_API_KEY)" ]; then \
		echo "Error: OPENAI_API_KEY not found in .env file"; \
		exit 1; \
	fi
	@CGO_ENABLED=1 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Build successful! Binary created: ./$(BINARY_NAME)"
	@echo "Signing binary..."
	@codesign --force --options runtime --timestamp -s - ./$(BINARY_NAME)
	@echo "Code signing successful!"

verify: ## Verify the code signature of the binary
	@echo "Verifying code signature..."
	@codesign -dvv ./$(BINARY_NAME)

verify-installed: ## Verify the code signature of the installed binary
	@echo "Verifying installed binary..."
	@codesign -dvv $(INSTALL_DIR)/$(BINARY_NAME)

install: build ## Build and install the binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"
	@echo "You can now run: $(BINARY_NAME) <filename>"

clean: ## Remove built binaries
	@echo "Cleaning up..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@$(GO) clean
	@echo "Clean complete!"

uninstall: ## Remove the installed binary from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstall complete!"

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v ./...

rebuild: clean build ## Clean and rebuild

reinstall: clean install ## Clean, rebuild, and reinstall

.DEFAULT_GOAL := help

