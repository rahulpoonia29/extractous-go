#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# All supported platforms
PLATFORMS=("linux_amd64" "darwin_amd64" "darwin_arm64" "windows_amd64")

echo "=============================================="
echo "  Building FFI for All Platforms"
echo "=============================================="
echo ""

# Track successes and failures
SUCCESS_COUNT=0
FAILED_PLATFORMS=()

for platform in "${PLATFORMS[@]}"; do
    echo ""
    echo "================================================"
    echo "  Building: $platform"
    echo "================================================"
    
    if "$SCRIPT_DIR/build_ffi.sh" "$platform"; then
        ((SUCCESS_COUNT++))
        echo "✓ $platform: SUCCESS"
    else
        FAILED_PLATFORMS+=("$platform")
        echo "✗ $platform: FAILED"
    fi
    
    echo ""
done

echo ""
echo "=============================================="
echo "  Build Summary"
echo "=============================================="
echo "Successful: $SUCCESS_COUNT / ${#PLATFORMS[@]}"

if [ ${#FAILED_PLATFORMS[@]} -gt 0 ]; then
    echo "Failed platforms:"
    for platform in "${FAILED_PLATFORMS[@]}"; do
        echo "  - $platform"
    done
    exit 1
else
    echo "All platforms built successfully!"
fi

echo ""
echo "Native libraries:"
for platform in "${PLATFORMS[@]}"; do
    native_dir="$ROOT_DIR/native/$platform"
    if [ -d "$native_dir" ]; then
        echo ""
        echo "[$platform]"
        ls -lh "$native_dir" | grep -E "\.(so|dylib|dll)$" | awk '{print "  " $9, "(" $5 ")"}'
    fi
done
echo ""
