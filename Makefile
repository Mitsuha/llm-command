# Makefile for command-llm

.PHONY: build install clean help

# Variables
OUTPUT_DIR := ./output
INSTALL_DIR := /usr/local/llm-command
CMD_DIR := ./cmd

# Default target
all: build

# Build all packages in cmd directory
build:
	@echo "🔨 Building packages..."
	@mkdir -p $(OUTPUT_DIR)
	@for dir in $(CMD_DIR)/*/; do \
		if [ -d "$$dir" ]; then \
			pkg_name=$$(basename "$$dir"); \
			echo "  Building $$pkg_name..."; \
			go build -o $(OUTPUT_DIR)/$$pkg_name $$dir; \
		fi; \
	done
	@echo "✅ Build completed. Binaries are in $(OUTPUT_DIR)/"

# Install binaries to system directory
install: build
	@echo "📦 Installing binaries to $(INSTALL_DIR)..."
	@sudo mkdir -p $(INSTALL_DIR)
	@sudo cp $(OUTPUT_DIR)/* $(INSTALL_DIR)/
	@echo "✅ Installation completed!"
	@echo ""
	@echo "📋 Setup Instructions:"
	@echo "===================="
	@echo "Add the following line to your shell profile:"
	@echo ""
	@echo "  For bash (~/.bashrc or ~/.bash_profile):"
	@echo "    export PATH=\"$(INSTALL_DIR):\$$PATH\""
	@echo ""
	@echo "  For zsh (~/.zshrc):"
	@echo "    export PATH=\"$(INSTALL_DIR):\$$PATH\""
	@echo ""
	@echo "  For fish (~/.config/fish/config.fish):"
	@echo "    set -gx PATH $(INSTALL_DIR) \$$PATH"
	@echo ""
	@echo "Then restart your terminal or run:"
	@echo "  source ~/.bashrc    # for bash"
	@echo "  source ~/.zshrc     # for zsh"
	@echo ""
	@echo "🎉 After setup, you can use: plz <description>"

# Clean output directory
clean:
	@echo "🧹 Cleaning output directory..."
	@rm -rf $(OUTPUT_DIR)
	@echo "✅ Clean completed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build   - Build all packages in ./cmd directory"
	@echo "  install - Install binaries to system (requires sudo)"
	@echo "  clean   - Clean the output directory"
	@echo "  help    - Show this help message"