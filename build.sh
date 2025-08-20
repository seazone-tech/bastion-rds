#!/bin/bash

# Build script for RDS Connector
set -e

APP_NAME="rds-connector"
VERSION="2.0"
BUILD_DIR="build"

echo "🏗️  RDS Connector Build Script v${VERSION}"
echo ""

# Clean build directory
echo "🧹 Limpando diretório de build..."
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Download dependencies
echo "📦 Baixando dependências..."
go mod download
go mod tidy

# Build for Linux (AMD64)
echo "🐧 Compilando para Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${APP_NAME}-linux-amd64 .

# Build for Windows (AMD64)  
echo "🪟 Compilando para Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${APP_NAME}-windows-amd64.exe .

# Build for macOS (AMD64)
echo "🍎 Compilando para macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${APP_NAME}-macos-amd64 .

# Build for macOS (ARM64 - Apple Silicon)
echo "🍎 Compilando para macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${APP_NAME}-macos-arm64 .

# Show results
echo ""
echo "✅ Build concluído com sucesso!"
echo ""
echo "📁 Executáveis criados em ${BUILD_DIR}/:"
ls -la ${BUILD_DIR}/

echo ""
echo "🚀 Para testar:"
echo "   Linux:   ./${BUILD_DIR}/${APP_NAME}-linux-amd64"
echo "   Windows: ./${BUILD_DIR}/${APP_NAME}-windows-amd64.exe"
echo "   macOS:   ./${BUILD_DIR}/${APP_NAME}-macos-amd64"
echo ""
echo "📦 Para criar pacotes de release: ./package.sh"