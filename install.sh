#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="wsoule/new-dotfiles"
BINARY_NAME="dotfiles"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    i386|i686) ARCH="386" ;;
    *) echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}" && exit 1 ;;
esac

case $OS in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo -e "${RED}❌ Unsupported OS: $OS${NC}" && exit 1 ;;
esac

echo -e "${BLUE}🛠  Dotfiles Manager Installer${NC}"
echo "=================================="
echo ""
echo "Installing for: $OS/$ARCH"
echo ""

# Check if running as root (not recommended)
if [[ $EUID -eq 0 ]]; then
    echo -e "${YELLOW}⚠️  Warning: Running as root. Consider running as regular user.${NC}"
fi

# Function to install from GitHub releases
install_from_releases() {
    echo -e "${BLUE}📦 Installing from GitHub releases...${NC}"

    # Get latest release info
    LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$LATEST_RELEASE" ]; then
        echo -e "${RED}❌ Could not fetch latest release info${NC}"
        return 1
    fi

    echo "Latest release: $LATEST_RELEASE"

    # Construct download URL
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"

    echo "Downloading: $DOWNLOAD_URL"

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Download and extract
    if ! curl -L -o "${BINARY_NAME}.tar.gz" "$DOWNLOAD_URL"; then
        echo -e "${RED}❌ Download failed${NC}"
        rm -rf "$TEMP_DIR"
        return 1
    fi

    tar -xzf "${BINARY_NAME}.tar.gz"

    # Install binary
    if [[ ! -w "$INSTALL_DIR" ]]; then
        echo -e "${YELLOW}🔐 Installing to $INSTALL_DIR (requires sudo)${NC}"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$BINARY_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Cleanup
    cd - > /dev/null
    rm -rf "$TEMP_DIR"

    echo -e "${GREEN}✅ Installed successfully!${NC}"
    return 0
}

# Function to build from source
install_from_source() {
    echo -e "${BLUE}🔨 Building from source...${NC}"

    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed. Please install Go first:${NC}"
        echo "   macOS: brew install go"
        echo "   Linux: Follow instructions at https://golang.org/doc/install"
        exit 1
    fi

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Clone repository
    echo "Cloning repository..."
    if ! git clone "https://github.com/$REPO.git" .; then
        echo -e "${RED}❌ Failed to clone repository${NC}"
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    # Build
    echo "Building..."
    if ! go build -o "$BINARY_NAME" .; then
        echo -e "${RED}❌ Build failed${NC}"
        rm -rf "$TEMP_DIR"
        exit 1
    fi

    # Install binary
    if [[ ! -w "$INSTALL_DIR" ]]; then
        echo -e "${YELLOW}🔐 Installing to $INSTALL_DIR (requires sudo)${NC}"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$BINARY_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Cleanup
    cd - > /dev/null
    rm -rf "$TEMP_DIR"

    echo -e "${GREEN}✅ Built and installed successfully!${NC}"
}

# Try installation methods
echo -e "${BLUE}Attempting installation...${NC}"

# Try GitHub releases first, fall back to source
if ! install_from_releases; then
    echo -e "${YELLOW}⚠️  Release installation failed, trying source build...${NC}"
    install_from_source
fi

# Verify installation
echo ""
echo -e "${BLUE}🧪 Verifying installation...${NC}"

if command -v "$BINARY_NAME" &> /dev/null; then
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo ""
    echo -e "${GREEN}🎉 Dotfiles Manager is ready!${NC}"
    echo ""
    echo "Get started with:"
    echo -e "  ${BLUE}dotfiles onboard${NC}          # Complete developer setup"
    echo -e "  ${BLUE}dotfiles init${NC}             # Initialize configuration"
    echo -e "  ${BLUE}dotfiles github setup${NC}     # Set up GitHub SSH"
    echo ""
    echo "For help:"
    echo -e "  ${BLUE}dotfiles --help${NC}"
    echo ""

    # Show version
    "$BINARY_NAME" --version 2>/dev/null || echo "Version: Latest"
else
    echo -e "${RED}❌ Installation verification failed${NC}"
    echo "The binary may not be in your PATH. Try:"
    echo "  $INSTALL_DIR/$BINARY_NAME --help"
    exit 1
fi
