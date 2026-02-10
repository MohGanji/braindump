#!/bin/bash
# braindump universal installer
# Usage: curl -fsSL https://raw.githubusercontent.com/MohGanji/braindump/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# GitHub repository
REPO="MohGanji/braindump"
GITHUB_RELEASES="https://github.com/${REPO}/releases/latest/download"

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin)
            echo "darwin"
            ;;
        Linux)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Print colored message
print_msg() {
    local color=$1
    shift
    echo -e "${color}$@${NC}"
}

# Print success message
success() {
    print_msg "$GREEN" "✓ $@"
}

# Print error message
error() {
    print_msg "$RED" "✗ $@"
}

# Print info message
info() {
    print_msg "$BLUE" "→ $@"
}

# Print warning message
warn() {
    print_msg "$YELLOW" "⚠ $@"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Determine install directory
get_install_dir() {
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
    elif [ -w "$HOME/.local/bin" ]; then
        mkdir -p "$HOME/.local/bin"
        echo "$HOME/.local/bin"
    else
        mkdir -p "$HOME/bin"
        echo "$HOME/bin"
    fi
}

# Update shell RC file
update_shell_rc() {
    local install_dir=$1
    local rc_file=""

    # Detect shell
    if [ -n "$BASH_VERSION" ]; then
        rc_file="$HOME/.bashrc"
        [ -f "$HOME/.bash_profile" ] && rc_file="$HOME/.bash_profile"
    elif [ -n "$ZSH_VERSION" ]; then
        rc_file="$HOME/.zshrc"
    else
        # Try to detect from SHELL variable
        case "$SHELL" in
            */bash)
                rc_file="$HOME/.bashrc"
                ;;
            */zsh)
                rc_file="$HOME/.zshrc"
                ;;
        esac
    fi

    if [ -n "$rc_file" ] && [ -f "$rc_file" ]; then
        # Check if already in PATH
        if ! echo "$PATH" | grep -q "$install_dir"; then
            info "Adding $install_dir to PATH in $rc_file"
            echo "" >> "$rc_file"
            echo "# Added by braindump installer" >> "$rc_file"
            echo "export PATH=\"$install_dir:\$PATH\"" >> "$rc_file"
            success "Updated $rc_file"
            return 0
        fi
    fi
    return 1
}

# Download pre-built binary
download_binary() {
    local os=$1
    local arch=$2
    local install_dir=$3

    local binary_name="braindump-${os}-${arch}"
    [ "$os" = "windows" ] && binary_name="${binary_name}.exe"

    local download_url="${GITHUB_RELEASES}/${binary_name}"
    local temp_file="/tmp/braindump-install-$$"

    info "Downloading pre-built binary from GitHub..."

    if command_exists curl; then
        if curl -fsSL "$download_url" -o "$temp_file" 2>/dev/null; then
            chmod +x "$temp_file"

            # Try to move to install dir
            if [ -w "$install_dir" ]; then
                mv "$temp_file" "$install_dir/braindump"
            else
                info "Need sudo to install to $install_dir"
                sudo mv "$temp_file" "$install_dir/braindump"
            fi

            success "Installed braindump to $install_dir/braindump"
            return 0
        fi
    elif command_exists wget; then
        if wget -q "$download_url" -O "$temp_file" 2>/dev/null; then
            chmod +x "$temp_file"

            if [ -w "$install_dir" ]; then
                mv "$temp_file" "$install_dir/braindump"
            else
                info "Need sudo to install to $install_dir"
                sudo mv "$temp_file" "$install_dir/braindump"
            fi

            success "Installed braindump to $install_dir/braindump"
            return 0
        fi
    fi

    # Download failed
    rm -f "$temp_file"
    return 1
}

# Build from source
build_from_source() {
    local install_dir=$1

    info "Pre-built binary not available, building from source..."

    if ! command_exists go; then
        error "Go is not installed and pre-built binary is not available"
        echo ""
        echo "Please either:"
        echo "  1. Install Go: https://golang.org/dl/"
        echo "  2. Wait for pre-built binaries to be published"
        exit 1
    fi

    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    info "Cloning repository..."
    if command_exists git; then
        git clone "https://github.com/${REPO}.git" braindump-src >/dev/null 2>&1
        cd braindump-src
    else
        error "Git is not installed"
        exit 1
    fi

    info "Building binary..."
    go build -o braindump . >/dev/null 2>&1

    if [ -w "$install_dir" ]; then
        mv braindump "$install_dir/braindump"
    else
        info "Need sudo to install to $install_dir"
        sudo mv braindump "$install_dir/braindump"
    fi

    # Cleanup
    cd /tmp
    rm -rf "$temp_dir"

    success "Built and installed braindump to $install_dir/braindump"
}

# Main installation
main() {
    echo ""
    print_msg "$BLUE" "╔════════════════════════════════════╗"
    print_msg "$BLUE" "║  braindump Universal Installer  ║"
    print_msg "$BLUE" "╚════════════════════════════════════╝"
    echo ""

    # Check if already installed
    if command_exists braindump; then
        warn "braindump is already installed at: $(which braindump)"
        echo ""

        # Skip prompt if non-interactive or NOTES_SKIP_PROMPT is set
        if [ -t 0 ] && [ -z "$NOTES_SKIP_PROMPT" ]; then
            read -p "Reinstall? (y/N) " -n 1 -r
            echo ""
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Installation cancelled"
                exit 0
            fi
        else
            info "Skipping reinstall (non-interactive mode or already installed)"
            success "braindump is ready to use!"
            exit 0
        fi
    fi

    # Detect system
    local os=$(detect_os)
    local arch=$(detect_arch)

    info "Detected system: $os/$arch"

    if [ "$os" = "unknown" ] || [ "$arch" = "unknown" ]; then
        error "Unsupported system: $os/$arch"
        exit 1
    fi

    # Determine install location
    local install_dir=$(get_install_dir)
    info "Install directory: $install_dir"
    echo ""

    # Try to download pre-built binary first
    if download_binary "$os" "$arch" "$install_dir"; then
        echo ""
        success "Pre-built binary installed successfully!"
    else
        # Fall back to building from source
        warn "Could not download pre-built binary"
        build_from_source "$install_dir"
    fi

    echo ""

    # Update PATH if needed
    local path_updated=0
    if ! echo "$PATH" | grep -q "$install_dir"; then
        if update_shell_rc "$install_dir"; then
            path_updated=1
        fi
    fi

    # Verify installation
    if [ -f "$install_dir/braindump" ]; then
        success "Installation complete!"
        echo ""

        # Show next steps
        print_msg "$GREEN" "╔════════════════════════════════════╗"
        print_msg "$GREEN" "║         Installation Complete!     ║"
        print_msg "$GREEN" "╚════════════════════════════════════╝"
        echo ""

        if [ $path_updated -eq 1 ]; then
            warn "PATH was updated. Run this command to use braindump now:"
            echo ""
            echo "    export PATH=\"$install_dir:\$PATH\""
            echo ""
            echo "Or start a new terminal session."
            echo ""
        fi

        echo "Try it out:"
        echo ""
        if [ $path_updated -eq 1 ]; then
            echo "    export PATH=\"$install_dir:\$PATH\""
        fi
        echo "    braindump add test --title \"Hello\" --content \"World\""
        echo "    braindump search \"hello\""
        echo "    braindump list"
        echo "    braindump help"
        echo ""

        # Try to add to PATH for current session
        export PATH="$install_dir:$PATH"

        if command_exists braindump; then
            success "braindump command is ready to use!"
            echo ""
            echo "Current version:"
            braindump help | head -n 1
        fi
    else
        error "Installation failed"
        exit 1
    fi
}

# Run main function
main "$@"
