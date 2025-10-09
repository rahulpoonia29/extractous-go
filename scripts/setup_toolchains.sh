#!/usr/bin/env bash
set -euo pipefail

echo "=============================================="
echo "  Rust Cross-Compilation Toolchain Setup"
echo "=============================================="
echo ""

# Check if rustup is installed
if ! command -v rustup &> /dev/null; then
    echo "Error: rustup is not installed."
    echo "Please install from: https://rustup.rs/"
    exit 1
fi

echo "Current Rust version:"
rustc --version
echo ""

# Define target triples for each platform
declare -A TARGETS=(
    ["linux_amd64"]="x86_64-unknown-linux-gnu"
    ["darwin_amd64"]="x86_64-apple-darwin"
    ["darwin_arm64"]="aarch64-apple-darwin"
    ["windows_amd64"]="x86_64-pc-windows-gnu"
)

echo "Installing cross-compilation targets..."
echo ""

for platform in "${!TARGETS[@]}"; do
    target="${TARGETS[$platform]}"
    echo "[$platform] Installing target: $target"
    
    if rustup target list --installed | grep -q "^$target$"; then
        echo "  ✓ Already installed"
    else
        rustup target add "$target"
        echo "  ✓ Installed successfully"
    fi
    echo ""
done

echo "=============================================="
echo "  Toolchain Setup Complete!"
echo "=============================================="
echo ""
echo "Installed targets:"
rustup target list --installed | grep -E "(linux-gnu|apple-darwin|windows-gnu)"
echo ""

# Platform-specific instructions
echo "=============================================="
echo "  Platform-Specific Requirements"
echo "=============================================="
echo ""

echo "Linux → macOS cross-compilation:"
echo "  Requires osxcross: https://github.com/tpoechtrager/osxcross"
echo ""

echo "Linux → Windows cross-compilation:"
echo "  Install mingw-w64:"
echo "    Ubuntu/Debian: sudo apt install mingw-w64"
echo "    Fedora/RHEL:   sudo dnf install mingw64-gcc"
echo ""

echo "macOS → Linux cross-compilation:"
echo "  Requires Docker or VM (cross-compilation toolchain)"
echo ""

echo "macOS → Windows cross-compilation:"
echo "  Install mingw-w64:"
echo "    brew install mingw-w64"
echo ""

echo "Note: Native compilation on each platform is recommended"
echo "      for production builds."
echo ""
