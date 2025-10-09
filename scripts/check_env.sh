#!/usr/bin/env bash
set -euo pipefail

echo "=============================================="
echo "  Environment Check"
echo "=============================================="
echo ""

ERRORS=0

# Check Rust
if command -v rustc &> /dev/null; then
    echo "✓ Rust:    $(rustc --version)"
    echo "  Cargo:   $(cargo --version)"
    echo "  Rustup:  $(rustup --version | head -n1)"
else
    echo "✗ Rust: Not installed"
    ((ERRORS++))
fi
echo ""

# Check build tools
echo "Build Tools:"
for tool in gcc make; do
    if command -v "$tool" &> /dev/null; then
        echo "  ✓ $tool"
    else
        echo "  ✗ $tool: Not found"
        ((ERRORS++))
    fi
done
echo ""

# Check Python tools (for wheel extraction)
echo "Python Tools (for wheel extraction):"
for tool in curl jq unzip sha256sum; do
    if command -v "$tool" &> /dev/null; then
        echo "  ✓ $tool"
    else
        echo "  ✗ $tool: Not found"
        ((ERRORS++))
    fi
done
echo ""

# Check cross-compilation targets
echo "Installed Rust Targets:"
rustup target list --installed | grep -E "(linux-gnu|apple-darwin|windows-gnu)" | sed 's/^/  /'
echo ""

# Check native libraries
echo "Native Libraries Status:"
PLATFORMS=("linux_amd64" "darwin_amd64" "darwin_arm64" "windows_amd64")
for platform in "${PLATFORMS[@]}"; do
    native_dir="native/$platform"
    if [ -d "$native_dir" ] && [ -n "$(ls -A "$native_dir" 2>/dev/null)" ]; then
        lib_count=$(find "$native_dir" -type f \( -name "*.so" -o -name "*.dylib" -o -name "*.dll" \) | wc -l)
        echo "  ✓ $platform: $lib_count libraries"
    else
        echo "  ✗ $platform: No libraries found"
    fi
done
echo ""

if [ $ERRORS -gt 0 ]; then
    echo "=============================================="
    echo "  Errors Found: $ERRORS"
    echo "=============================================="
    exit 1
else
    echo "=============================================="
    echo "  Environment OK!"
    echo "=============================================="
fi
