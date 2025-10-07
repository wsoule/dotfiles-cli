#!/bin/bash
# Dotfiles Manager - Linux Installation Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}🚀 Dotfiles Manager - Linux Installer${NC}"
echo "Architecture: $ARCH"
echo ""

# Get latest release
LATEST_RELEASE=$(curl -s https://api.github.com/repos/wsoule/dotfiles-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}Failed to get latest release${NC}"
    exit 1
fi

echo "Latest version: $LATEST_RELEASE"
echo ""

# Download URL
DOWNLOAD_URL="https://github.com/wsoule/dotfiles-cli/releases/download/${LATEST_RELEASE}/dotfiles_Linux_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"
echo ""

# Create temp directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download
echo -e "${YELLOW}Downloading...${NC}"
if ! curl -L "$DOWNLOAD_URL" -o dotfiles.tar.gz; then
    echo -e "${RED}Download failed${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Extract
echo -e "${YELLOW}Extracting...${NC}"
tar -xzf dotfiles.tar.gz

# Make executable
chmod +x dotfiles

# Install
INSTALL_DIR="/usr/local/bin"

if [ -w "$INSTALL_DIR" ]; then
    # Can write without sudo
    echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
    mv dotfiles "$INSTALL_DIR/"
else
    # Need sudo
    echo -e "${YELLOW}Installing to $INSTALL_DIR (requires sudo)...${NC}"
    sudo mv dotfiles "$INSTALL_DIR/"
fi

# Cleanup
cd -
rm -rf "$TMP_DIR"

# Verify installation
if command -v dotfiles &> /dev/null; then
    echo ""
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo ""
    echo "Dotfiles Manager is now installed. Try:"
    echo "  dotfiles --help"
    echo "  dotfiles onboard"
    echo ""
else
    echo -e "${RED}Installation failed. Binary not found in PATH.${NC}"
    exit 1
fi
