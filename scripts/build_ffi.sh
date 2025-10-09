#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
FFI_DIR="$ROOT_DIR/ffi"

# Platform mapping
declare -A RUST_TARGETS=(
    ["linux_amd64"]="x86_64-unknown-linux-gnu"
    ["darwin_amd64"]="x86_64-apple-darwin"
    ["darwin_arm64"]="aarch64-apple-darwin"
    ["windows_amd64"]="x86_64-pc-windows-gnu"
)

declare -A LIB_EXTENSIONS=(
    ["linux_amd64"]="so"
    ["darwin_amd64"]="dylib"
    ["darwin_arm64"]="dylib"
    ["windows_amd64"]="dll"
)

declare -A LIB_PREFIXES=(
    ["linux_amd64"]="lib"
    ["darwin_amd64"]="lib"
    ["darwin_arm64"]="lib"
    ["windows_amd64"]=""
)

# Get target platform (from argument or detect current)
if [ -n "${1:-}" ]; then
    PLATFORM="$1"
else
    PLATFORM=$("$SCRIPT_DIR/detect_platform.sh")
fi

# Validate platform
if [ -z "${RUST_TARGETS[$PLATFORM]:-}" ]; then
    echo "Error: Unsupported platform: $PLATFORM"
    echo "Supported: ${!RUST_TARGETS[@]}"
    exit 1
fi

RUST_TARGET="${RUST_TARGETS[$PLATFORM]}"
LIB_EXT="${LIB_EXTENSIONS[$PLATFORM]}"
LIB_PREFIX="${LIB_PREFIXES[$PLATFORM]}"

# Build profile
PROFILE="${CARGO_PROFILE:-release}"
BUILD_FLAGS=()
if [ "$PROFILE" = "release" ]; then
    BUILD_FLAGS+=("--release")
fi

echo "=============================================="
echo "  Building extractous-ffi"
echo "=============================================="
echo "Platform:     $PLATFORM"
echo "Rust Target:  $RUST_TARGET"
echo "Profile:      $PROFILE"
echo ""

# Check if native libraries exist
NATIVE_DIR="$ROOT_DIR/native/$PLATFORM"
if [ ! -d "$NATIVE_DIR" ] || [ -z "$(ls -A "$NATIVE_DIR" 2>/dev/null)" ]; then
    echo "ERROR: Native libraries not found for $PLATFORM"
    echo "Expected at: $NATIVE_DIR"
    echo ""
    echo "Please run: make extract-wheels"
    exit 1
fi

echo "Native libraries found:"
ls -lh "$NATIVE_DIR" | tail -n +2
echo ""

# Build FFI library
cd "$FFI_DIR"
echo "Building Rust FFI library..."
echo "Command: cargo build ${BUILD_FLAGS[*]} --target $RUST_TARGET"
echo ""

cargo build "${BUILD_FLAGS[@]}" --target "$RUST_TARGET"

# Locate built library
TARGET_DIR="$FFI_DIR/target/$RUST_TARGET/$PROFILE"
LIB_NAME="${LIB_PREFIX}extractous_ffi"
if [ "$PLATFORM" = "windows_amd64" ]; then
    LIB_FILE="$TARGET_DIR/extractous_ffi.dll"
else
    LIB_FILE="$TARGET_DIR/${LIB_NAME}.${LIB_EXT}"
fi

if [ ! -f "$LIB_FILE" ]; then
    echo "ERROR: Built library not found at: $LIB_FILE"
    exit 1
fi

echo ""
echo "Built library: $LIB_FILE"
echo "Size: $(du -h "$LIB_FILE" | cut -f1)"
echo ""

# Copy to native directory
echo "Copying to: $NATIVE_DIR"
mkdir -p "$NATIVE_DIR"
cp -v "$LIB_FILE" "$NATIVE_DIR/"
echo ""

# Verify RPATH (Linux/macOS only)
if [ "$PLATFORM" = "linux_amd64" ]; then
    echo "Verifying RPATH..."
    if command -v readelf &> /dev/null; then
        readelf -d "$LIB_FILE" | grep -E "(RPATH|RUNPATH)" || echo "  (no RPATH/RUNPATH)"
    fi
elif [[ "$PLATFORM" == darwin_* ]]; then
    echo "Verifying LC_RPATH..."
    if command -v otool &> /dev/null; then
        otool -l "$LIB_FILE" | grep -A 2 "LC_RPATH" || echo "  (no LC_RPATH)"
    fi
fi

echo ""
echo "=============================================="
echo "  Build Complete!"
echo "=============================================="
echo "Platform:  $PLATFORM"
echo "Library:   $(basename "$LIB_FILE")"
echo "Location:  $NATIVE_DIR"
echo ""
