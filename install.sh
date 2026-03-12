#!/bin/sh
set -e

REPO="getopsbudget/opsbudget-cli"
BINARY="opsbudget"
INSTALL_DIR="/usr/local/bin"

main() {
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "$os" in
        linux)  os="linux" ;;
        darwin) os="darwin" ;;
        *)
            echo "Error: unsupported OS: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64)  arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *)
            echo "Error: unsupported architecture: $arch"
            exit 1
            ;;
    esac

    # Get latest release tag
    latest=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$latest" ]; then
        echo "Error: could not determine latest release"
        exit 1
    fi

    url="https://github.com/${REPO}/releases/download/v${latest}/opsbudget_${latest}_${os}_${arch}.tar.gz"
    echo "Downloading opsbudget v${latest} for ${os}/${arch}..."

    tmpdir=$(mktemp -d)
    trap 'rm -rf "$tmpdir"' EXIT

    curl -sSL "$url" | tar -xz -C "$tmpdir"

    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmpdir/$BINARY" "$INSTALL_DIR/$BINARY"
    else
        echo "Installing to $INSTALL_DIR (requires sudo)..."
        sudo mv "$tmpdir/$BINARY" "$INSTALL_DIR/$BINARY"
    fi
    chmod +x "$INSTALL_DIR/$BINARY"

    echo ""
    echo "✓ opsbudget installed successfully (v${latest})"
    echo ""
    echo "Get started:"
    echo "  opsbudget login"
    echo "  opsbudget ping add https://yoursite.com"
    echo ""
    echo "Docs: https://opsbudget.com/docs"
}

main
