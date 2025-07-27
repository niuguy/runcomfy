#!/bin/bash

# RunComfy Installation Script for RunPod
set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="runcomfy"
VERSION="latest"
GITHUB_REPO="your-username/runcomfy"  # Update this with your actual repo

echo "🚀 Installing RunComfy..."

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Download URL (adjust based on your release naming)
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/runcomfy-${OS}-${ARCH}"

echo "📥 Downloading RunComfy for ${OS}/${ARCH}..."

# Download binary
if command -v curl >/dev/null 2>&1; then
    curl -L -o "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"
elif command -v wget >/dev/null 2>&1; then
    wget -O "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"
else
    echo "❌ Neither curl nor wget is available. Please install one of them."
    exit 1
fi

# Make executable
chmod +x "/tmp/${BINARY_NAME}"

# Move to install directory
sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

echo "✅ RunComfy installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
echo ""
echo "🔧 Usage:"
echo "  runcomfy analyze workflow.json"
echo "  runcomfy scan"
echo "  runcomfy install workflow.json"
echo ""
echo "📖 Run 'runcomfy --help' for more information"