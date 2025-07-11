#!/bin/sh

set -e

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *)
    echo "Error: Unsupported operating system: $OS"
    echo "Supported: linux, darwin (macOS)"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Error: Unsupported architecture: $ARCH"
    echo "Supported: x86_64 (amd64), aarch64/arm64"
    exit 1
    ;;
esac

# Construct download URL
BINARY_NAME="compartment-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/lcmen/compartment/releases/latest/download/${BINARY_NAME}"

echo "Detecting platform: ${OS}-${ARCH}"
echo "Downloading from: ${DOWNLOAD_URL}"

# Download the binary
if command -v curl >/dev/null 2>&1; then
  curl -L "${DOWNLOAD_URL}" -o compartment
elif command -v wget >/dev/null 2>&1; then
  wget "${DOWNLOAD_URL}" -O compartment
else
  echo "Error: Neither curl nor wget is available"
  echo "Please install curl or wget and try again"
  exit 1
fi

# Make executable
chmod +x compartment

echo "✓ compartment downloaded successfully"
echo "Run './compartment --help' to get started"