#!/bin/bash

# Package script for RDS Connector releases
set -e

APP_NAME="rds-connector"
VERSION="2.0"
BUILD_DIR="build"
PACKAGE_DIR="packages"

echo "📦 RDS Connector Package Script v${VERSION}"
echo ""

# Check if builds exist
if [ ! -d "${BUILD_DIR}" ]; then
    echo "❌ Diretório de build não encontrado. Execute ./build.sh primeiro."
    exit 1
fi

# Create package directory
echo "📁 Criando diretório de pacotes..."
rm -rf ${PACKAGE_DIR}
mkdir -p ${PACKAGE_DIR}

echo "🗜️  Criando pacotes comprimidos..."

# Create Linux package
if [ -f "${BUILD_DIR}/${APP_NAME}-linux-amd64" ]; then
    echo "  📦 Linux AMD64..."
    cd ${BUILD_DIR}
    tar -czf ../${PACKAGE_DIR}/${APP_NAME}-linux-amd64-v${VERSION}.tar.gz ${APP_NAME}-linux-amd64
    cd ..
fi

# Create Windows package
if [ -f "${BUILD_DIR}/${APP_NAME}-windows-amd64.exe" ]; then
    echo "  📦 Windows AMD64..."
    cd ${BUILD_DIR}
    zip -q ../${PACKAGE_DIR}/${APP_NAME}-windows-amd64-v${VERSION}.zip ${APP_NAME}-windows-amd64.exe
    cd ..
fi

# Create macOS Intel package
if [ -f "${BUILD_DIR}/${APP_NAME}-macos-amd64" ]; then
    echo "  📦 macOS Intel..."
    cd ${BUILD_DIR}
    tar -czf ../${PACKAGE_DIR}/${APP_NAME}-macos-amd64-v${VERSION}.tar.gz ${APP_NAME}-macos-amd64
    cd ..
fi

# Create macOS Apple Silicon package
if [ -f "${BUILD_DIR}/${APP_NAME}-macos-arm64" ]; then
    echo "  📦 macOS Apple Silicon..."
    cd ${BUILD_DIR}
    tar -czf ../${PACKAGE_DIR}/${APP_NAME}-macos-arm64-v${VERSION}.tar.gz ${APP_NAME}-macos-arm64
    cd ..
fi

echo ""
echo "✅ Pacotes criados com sucesso!"
echo ""
echo "📁 Pacotes disponíveis em ${PACKAGE_DIR}/:"
ls -la ${PACKAGE_DIR}/

echo ""
echo "🔍 Tamanhos dos pacotes:"
du -h ${PACKAGE_DIR}/*

echo ""
echo "🚀 Instruções para distribuição:"
echo ""
echo "Linux/macOS:"
echo "  1. Extrair: tar -xzf ${APP_NAME}-linux-amd64-v${VERSION}.tar.gz"
echo "  2. Executar: chmod +x ${APP_NAME}-linux-amd64 && ./${APP_NAME}-linux-amd64"
echo ""
echo "Windows:"
echo "  1. Extrair: unzip ${APP_NAME}-windows-amd64-v${VERSION}.zip"
echo "  2. Executar: ${APP_NAME}-windows-amd64.exe"
echo ""
echo "📋 Requisitos para desenvolvedores:"
echo "  - kubectl instalado e configurado"
echo "  - Acesso ao cluster EKS com credenciais AWS"
echo "  - Permissões para criar pods nos namespaces stg-apps/prd-apps"