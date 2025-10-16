.PHONY: build install clean test verify help

# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    BINARY_EXT := .exe
    INSTALL_DIR := $(USERPROFILE)/bin
else
    BINARY_EXT :=
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        DETECTED_OS := Linux
        INSTALL_DIR := /usr/local/bin
    endif
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macOS
        INSTALL_DIR := /usr/local/bin
    endif
endif

# Variables
BINARY_NAME=revyu$(BINARY_EXT)
BUILD_DIR=.
GO=go
LDFLAGS=-ldflags "-X main.buildTimeAPIKey=$(OPENAI_API_KEY)"

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Detected OS: $(DETECTED_OS)"
	@echo "Install directory: $(INSTALL_DIR)"
	@echo ""
	@echo "Targets:"
	@echo "  \033[36mbuild\033[0m           Build the binary"
	@echo "  \033[36minstall\033[0m         Build and install the binary"
	@echo "  \033[36mverify\033[0m          Verify the code signature (macOS only)"
	@echo "  \033[36mclean\033[0m           Remove built binaries"
	@echo "  \033[36muninstall\033[0m       Remove the installed binary"
	@echo "  \033[36mrebuild\033[0m         Clean and rebuild"
	@echo "  \033[36mreinstall\033[0m       Clean, rebuild, and reinstall"
	@echo "  \033[36mtest\033[0m            Run tests"
	@echo "  \033[36mhelp\033[0m            Show this help message"

build: ## Build the binary
	@echo "Building $(BINARY_NAME) for $(DETECTED_OS)..."
	@if [ -z "$(OPENAI_API_KEY)" ]; then \
		echo "Warning: OPENAI_API_KEY not found - building without embedded key"; \
		echo "You'll need to set OPENAI_API_KEY environment variable or create .env file"; \
		$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME); \
	else \
		echo "Building with embedded OpenAI API key..."; \
		$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME); \
	fi
	@echo "Build successful! Binary created: ./$(BINARY_NAME)"
ifeq ($(DETECTED_OS),macOS)
	@echo "Code signing binary for macOS..."
	@codesign --force --options runtime --timestamp -s - ./$(BINARY_NAME) 2>/dev/null || echo "Warning: Code signing failed (optional)"
endif

verify: ## Verify the code signature (macOS only)
ifeq ($(DETECTED_OS),macOS)
	@echo "Verifying code signature..."
	@codesign -dvv ./$(BINARY_NAME)
else
	@echo "Code signing verification is only available on macOS"
endif

install: build ## Build and install the binary
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
ifeq ($(DETECTED_OS),Windows)
	@if not exist "$(INSTALL_DIR)" mkdir "$(INSTALL_DIR)"
	@copy /Y $(BUILD_DIR)\$(BINARY_NAME) $(INSTALL_DIR)\$(BINARY_NAME)
	@echo "Installation complete!"
	@echo "Make sure $(INSTALL_DIR) is in your PATH"
	@echo "You can now run: $(BINARY_NAME) <filename>"
else
	@mkdir -p $(INSTALL_DIR)
	@if [ -w $(INSTALL_DIR) ]; then \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME); \
		chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME); \
		sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "Installation complete!"
	@echo "You can now run: revyu <filename>"
endif

clean: ## Remove built binaries
	@echo "Cleaning up..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@$(GO) clean
	@echo "Clean complete!"

uninstall: ## Remove the installed binary
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
ifeq ($(DETECTED_OS),Windows)
	@del /F /Q $(INSTALL_DIR)\$(BINARY_NAME) 2>nul || echo "Binary not found or already removed"
else
	@if [ -w $(INSTALL_DIR)/$(BINARY_NAME) ]; then \
		rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
endif
	@echo "Uninstall complete!"

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v ./...

rebuild: clean build ## Clean and rebuild

reinstall: clean install ## Clean, rebuild, and reinstall

.DEFAULT_GOAL := help

