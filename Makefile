# Go parameters
BINARY_NAME=FTransferor
MAIN_PATH=.
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get
GOMOD=$(GO) mod
GOFMT=$(GO) fmt

# Build flags
LDFLAGS=-ldflags "-s -w"

# Output directories
BIN_DIR=bin

# Build targets for different platforms
.PHONY: all build clean test fmt check build-all windows linux macos

all: clean build

# Build for current platform
build:
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for all platforms
build-all: windows linux macos

# Build for Windows
windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)

# Build for Linux
linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-Linux $(MAIN_PATH)

# Build for MacOS
macos:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-MacOS $(MAIN_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)*

# Run tests
test:
	$(GOTEST) -v ./...

# Format code
fmt:
	$(GOFMT) ./...

# Check and tidy dependencies
check:
	$(GOMOD) tidy
	$(GOMOD) verify

# Install dependencies
deps:
	$(GOMOD) download

# Help target
help:
	@echo "Available targets:"
	@echo "  make          : Build for current platform"
	@echo "  make build    : Same as above"
	@echo "  make build-all: Build for all platforms"
	@echo "  make windows  : Build for Windows"
	@echo "  make linux    : Build for Linux"
	@echo "  make macos   : Build for MacOS"
	@echo "  make clean    : Clean build artifacts"
	@echo "  make test     : Run tests"
	@echo "  make fmt      : Format code"
	@echo "  make check    : Check and tidy dependencies"
	@echo "  make deps     : Install dependencies"
