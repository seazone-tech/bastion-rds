# RDS Connector - Makefile for cross-platform builds

APP_NAME=rds-connector
VERSION=2.0
BUILD_DIR=build

# Go build flags
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

# Default target
.PHONY: all
all: clean build-all

# Clean build directory
.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}
	mkdir -p ${BUILD_DIR}

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-macos

# Build for Linux (AMD64)
.PHONY: build-linux
build-linux:
	@echo "Building for Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME}-linux-amd64 .

# Build for Linux (ARM64)
.PHONY: build-linux-arm
build-linux-arm:
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME}-linux-arm64 .

# Build for Windows (AMD64)
.PHONY: build-windows
build-windows:
	@echo "Building for Windows AMD64..."
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe .

# Build for macOS (AMD64)
.PHONY: build-macos
build-macos:
	@echo "Building for macOS AMD64..."
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME}-macos-amd64 .

# Build for macOS (ARM64 - Apple Silicon)
.PHONY: build-macos-arm
build-macos-arm:
	@echo "Building for macOS ARM64..."
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME}-macos-arm64 .

# Build for current platform
.PHONY: build
build:
	@echo "Building for current platform..."
	go build ${LDFLAGS} -o ${BUILD_DIR}/${APP_NAME} .

# Run the application
.PHONY: run
run:
	go run .

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Test the application
.PHONY: test
test:
	go test ./...

# Install dependencies and build
.PHONY: install
install: deps build

# Development build (with debugging symbols)
.PHONY: dev
dev:
	go build -o ${BUILD_DIR}/${APP_NAME}-dev .

# Create release packages
.PHONY: package
package: build-all
	@echo "Creating release packages..."
	cd ${BUILD_DIR} && \
	tar -czf ${APP_NAME}-linux-amd64-v${VERSION}.tar.gz ${APP_NAME}-linux-amd64 && \
	tar -czf ${APP_NAME}-linux-arm64-v${VERSION}.tar.gz ${APP_NAME}-linux-arm64 && \
	zip ${APP_NAME}-windows-amd64-v${VERSION}.zip ${APP_NAME}-windows-amd64.exe && \
	tar -czf ${APP_NAME}-macos-amd64-v${VERSION}.tar.gz ${APP_NAME}-macos-amd64 && \
	tar -czf ${APP_NAME}-macos-arm64-v${VERSION}.tar.gz ${APP_NAME}-macos-arm64
	@echo "Release packages created in ${BUILD_DIR}/"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build for all platforms"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-linux  - Build for Linux AMD64"
	@echo "  build-windows- Build for Windows AMD64"
	@echo "  build-macos  - Build for macOS AMD64"
	@echo "  clean        - Clean build directory"
	@echo "  deps         - Download dependencies"
	@echo "  dev          - Development build"
	@echo "  install      - Install deps and build"
	@echo "  package      - Create release packages"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  help         - Show this help"