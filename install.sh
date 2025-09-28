#!/bin/bash

# Dotfiles Installation Script
# This script downloads and installs the latest version of dotfiles

set -e

# Configuration
REPO="wsoule/new-dotfiles"  # Change this to your GitHub username/repo
BINARY_NAME="dotfiles"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch

    case "$(uname -s)" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        MINGW* | MSYS* | CYGWIN*)
            os="windows"
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64 | amd64)
            arch="x86_64"
            ;;
        arm64 | aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local api_url="https://api.github.com/repos/${REPO}/releases/latest"

    if command -v curl >/dev/null 2>&1; then
        curl -s "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Download and install
install_dotfiles() {
    local platform
    local version
    local download_url
    local temp_dir
    local archive_name

    platform=$(detect_platform)
    version=$(get_latest_version)

    if [ -z "$version" ]; then
        print_error "Could not determine latest version"
        exit 1
    fi

    print_info "Latest version: ${version}"
    print_info "Platform: ${platform}"

    # Construct download URL
    archive_name="${BINARY_NAME}-${version}-${platform}"
    if [[ "$platform" == *"windows"* ]]; then
        archive_name="${archive_name}.zip"
    else
        archive_name="${archive_name}.tar.gz"
    fi

    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    print_info "Downloading from: ${download_url}"

    # Create temporary directory
    temp_dir=$(mktemp -d)
    cd "$temp_dir"

    # Download archive
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$archive_name" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$archive_name" "$download_url"
    else
        print_error "Neither curl nor wget is available"
        exit 1
    fi

    # Download and verify checksums for security
    local checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
    local checksums_file="checksums.txt"

    print_info "Verifying download integrity..."
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$checksums_file" "$checksums_url" 2>/dev/null || true
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$checksums_file" "$checksums_url" 2>/dev/null || true
    fi

    # Verify checksum if available and shasum is present
    if [ -f "$checksums_file" ] && command -v shasum >/dev/null 2>&1; then
        if shasum -a 256 -c "$checksums_file" --ignore-missing --quiet 2>/dev/null; then
            print_success "Download integrity verified"
        else
            print_warning "Could not verify download integrity (checksum mismatch or not found)"
            read -p "Continue anyway? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                print_error "Installation cancelled for security"
                exit 1
            fi
        fi
    else
        print_warning "Could not verify download integrity (checksums unavailable)"
    fi

    # Extract archive
    if [[ "$archive_name" == *.zip ]]; then
        unzip -q "$archive_name"
    else
        tar -xzf "$archive_name"
    fi

    # Find binary
    binary_path=""
    if [ -f "$BINARY_NAME" ]; then
        binary_path="$BINARY_NAME"
    elif [ -f "*/bin/$BINARY_NAME" ]; then
        binary_path="*/bin/$BINARY_NAME"
    else
        print_error "Binary not found in archive"
        exit 1
    fi

    # Install binary
    if [ -w "$INSTALL_DIR" ]; then
        cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        print_info "Installing to $INSTALL_DIR (requires sudo)"
        sudo cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Cleanup
    cd /
    rm -rf "$temp_dir"

    print_success "Successfully installed $BINARY_NAME to $INSTALL_DIR"
    print_info "Run '$BINARY_NAME --help' to get started"
}

# Main execution
main() {
    echo "🛠  Dotfiles Installer"
    echo "====================="
    echo

    # Check if already installed
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version
        current_version=$("$BINARY_NAME" --version 2>/dev/null | grep -o 'v[0-9.]*' || echo "unknown")
        print_warning "$BINARY_NAME is already installed (version: $current_version)"

        read -p "Do you want to update it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled"
            exit 0
        fi
    fi

    # Install
    install_dotfiles

    echo
    print_success "Installation complete!"
    echo
    print_info "Next steps:"
    echo "  1. Run: $BINARY_NAME setup"
    echo "  2. Run: $BINARY_NAME install"
    echo
}

# Run main function
main "$@"
