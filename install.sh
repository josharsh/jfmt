#!/bin/bash
set -e

# jfmt installer
# Usage: curl -fsSL https://raw.githubusercontent.com/josharsh/jfmt/main/install.sh | bash

VERSION="${JFMT_VERSION:-latest}"
INSTALL_DIR="${JFMT_INSTALL_DIR:-/usr/local/bin}"
REPO="josharsh/jfmt"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest version if not specified
if [ "$VERSION" = "latest" ]; then
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
fi

BINARY="jfmt-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY="${BINARY}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"

echo "Installing jfmt ${VERSION} for ${OS}/${ARCH}..."

# Download and install
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/jfmt"
chmod +x "${TMP_DIR}/jfmt"

# Install (may require sudo)
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_DIR}/jfmt" "${INSTALL_DIR}/jfmt"
else
    echo "Installing to ${INSTALL_DIR} requires sudo..."
    sudo mv "${TMP_DIR}/jfmt" "${INSTALL_DIR}/jfmt"
fi

echo "Installed jfmt to ${INSTALL_DIR}/jfmt"
echo ""
echo "Run 'jfmt -h' to get started!"
