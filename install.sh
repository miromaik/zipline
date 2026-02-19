#!/bin/bash
set -e

echo "Installing Zipline..."

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)
        case "$ARCH" in
            x86_64) PLATFORM="linux-amd64" ;;
            aarch64|arm64) PLATFORM="linux-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    Darwin*)
        case "$ARCH" in
            x86_64) PLATFORM="darwin-amd64" ;;
            arm64) PLATFORM="darwin-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    MINGW*|MSYS*|CYGWIN*)
        PLATFORM="windows-amd64"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="zipline"

if [[ "$PLATFORM" == "windows-amd64" ]]; then
    BINARY_NAME="zipline.exe"
fi

echo "Platform: $PLATFORM"
echo "Install directory: $INSTALL_DIR"

if [ -f "./zipline" ]; then
    echo "Installing from local build..."
    sudo cp ./zipline "$INSTALL_DIR/zipline"
    sudo chmod +x "$INSTALL_DIR/zipline"
elif command -v go &> /dev/null; then
    echo "Building from source..."
    go build -o "$BINARY_NAME"
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Error: Neither pre-built binary nor Go compiler found."
    echo "Please install Go or download a release from GitHub."
    exit 1
fi

echo "✓ Zipline installed successfully!"
echo ""
echo "Usage:"
echo "  zipline send file.pdf"
echo "  zipline get 123456"
echo ""
