#!/usr/bin/env bash
set -e

# Collect all native libraries for distribution
# Usage: ./scripts/collect-libs.sh <platform> <target> <lib_ext>

PLATFORM=$1  # e.g., linux_amd64
TARGET=$2    # e.g., x86_64-unknown-linux-gnu
LIB_EXT=$3   # e.g., so, dll, dylib

echo "=== Collecting Libraries ==="
echo "Platform: $PLATFORM | Target: $TARGET | Extension: $LIB_EXT"

# Create distribution structure
DIST_DIR="dist/$PLATFORM"
mkdir -p "$DIST_DIR/lib"
mkdir -p "$DIST_DIR/include"

# Determine OS-specific naming
if [[ "$PLATFORM" == *"windows"* ]]; then
    OS="Windows"
    MAIN_LIB_PREFIX=""
    TIKA_LIB_PREFIX=""
elif [[ "$PLATFORM" == *"darwin"* ]]; then
    OS="macOS"
    MAIN_LIB_PREFIX="lib"
    TIKA_LIB_PREFIX="lib"
else
    OS="Linux"
    MAIN_LIB_PREFIX="lib"
    TIKA_LIB_PREFIX="lib"
fi

# 1. Find and copy main FFI library
RELEASE_DIR="ffi/target/$TARGET/release"
MAIN_LIB="$RELEASE_DIR/${MAIN_LIB_PREFIX}extractous_ffi.$LIB_EXT"

if [ ! -f "$MAIN_LIB" ]; then
    echo "✗ Error: Main FFI library not found: $MAIN_LIB"
    exit 1
fi

echo "✓ Found main FFI library: $MAIN_LIB"
cp "$MAIN_LIB" "$DIST_DIR/lib/"

# 2. Find extractous build output directory
# Look for the canonical libs directory created by extractous build.rs
echo ""
echo "Searching for extractous dependencies..."

BUILD_BASE="$RELEASE_DIR/build"

# Find all extractous-*/out/libs directories, sort by modification time (newest first)
LIBS_DIR=$(find "$BUILD_BASE" -type d -path "*/extractous-*/out/libs" -printf "%T@ %p\n" 2>/dev/null | \
           sort -rn | \
           head -1 | \
           cut -d' ' -f2)

# Fallback for macOS (no -printf support)
if [ -z "$LIBS_DIR" ]; then
    LIBS_DIR=$(find "$BUILD_BASE" -type d -path "*/extractous-*/out/libs" -print 2>/dev/null | \
               xargs -0 stat -f "%m %N" 2>/dev/null | \
               sort -rn | \
               head -1 | \
               cut -d' ' -f2)
fi

if [ -z "$LIBS_DIR" ]; then
    echo "✗ Error: Could not find extractous out/libs directory"
    echo "Searched in: $BUILD_BASE/extractous-*/out/libs"
    exit 1
fi

# 3. Verify libtika_native exists
# Try both with and without prefix for Windows compatibility
TIKA_LIB="$LIBS_DIR/${TIKA_LIB_PREFIX}tika_native.$LIB_EXT"
if [ ! -f "$TIKA_LIB" ] && [ "$OS" = "Windows" ]; then
    # Try with lib prefix on Windows as fallback
    TIKA_LIB_ALT="$LIBS_DIR/libtika_native.$LIB_EXT"
    if [ -f "$TIKA_LIB_ALT" ]; then
        TIKA_LIB="$TIKA_LIB_ALT"
    fi
fi

if [ ! -f "$TIKA_LIB" ]; then
    echo "✗ Error: tika_native.$LIB_EXT not found in $LIBS_DIR"
    echo "Directory contents:"
    ls -lh "$LIBS_DIR" || echo "Directory not accessible"
    
    # Show all DLL/SO/DYLIB files to help debug
    echo ""
    echo "All native libraries found:"
    find "$LIBS_DIR" -name "*.$LIB_EXT" -o -name "*.dll" -o -name "*.so" -o -name "*.dylib" 2>/dev/null || echo "None found"
    exit 1
fi

echo "✓ Found libtika_native: $TIKA_LIB"

# 4. Copy ALL libraries from out/libs/
# These are all required dependencies bundled by GraalVM
echo ""
echo "Copying all native dependencies..."
cp "$LIBS_DIR"/*.$LIB_EXT "$DIST_DIR/lib/" || {
    echo "✗ Error: Failed to copy libraries"
    exit 1
}

# Count copied libraries
LIB_COUNT=$(find "$DIST_DIR/lib" -name "*.$LIB_EXT" | wc -l)
echo "✓ Copied $LIB_COUNT libraries"

# 5. Copy C header
HEADER="../include/extractous.h"
if [ -f "$HEADER" ]; then
    cp "$HEADER" "$DIST_DIR/include/"
    echo "✓ Copied C header"
fi

# 6. Display distribution contents
echo ""
echo "=== Distribution Contents ==="
echo "Libraries ($LIB_COUNT total):"
ls -lh "$DIST_DIR/lib/" | tail -n +2

if [ -d "$DIST_DIR/include" ] && [ -n "$(ls -A "$DIST_DIR/include" 2>/dev/null)" ]; then
    echo ""
    echo "Headers:"
    ls -lh "$DIST_DIR/include/"
fi

# 7. Verify dependencies (platform-specific)
echo ""
echo "=== Dependency Verification ==="

if [ "$OS" = "Linux" ]; then
    echo "Checking RPATH configuration..."
    for lib in "$DIST_DIR/lib/"*.$LIB_EXT; do
        LIB_NAME=$(basename "$lib")
        echo "  $LIB_NAME:"
        readelf -d "$lib" | grep -E "RPATH|RUNPATH" | sed 's/^/    /' || echo "    No RPATH set"
    done
    
    echo ""
    echo "Checking main FFI dependencies..."
    ldd "$DIST_DIR/lib/${MAIN_LIB_PREFIX}extractous_ffi.$LIB_EXT" || true

elif [ "$OS" = "macOS" ]; then
    echo "Checking install names..."
    for lib in "$DIST_DIR/lib/"*.$LIB_EXT; do
        LIB_NAME=$(basename "$lib")
        echo "  $LIB_NAME:"
        otool -L "$lib" | grep -v "$LIB_NAME:" | sed 's/^/    /'
    done

elif [ "$OS" = "Windows" ]; then
    echo "Windows DLL validation..."
    file "$DIST_DIR/lib"/*.dll 2>/dev/null || echo "  DLLs present"
fi

# 8. Calculate total size
echo ""
echo "=== Distribution Size ==="
echo "Total library size: $(du -sh "$DIST_DIR/lib" | cut -f1)"
echo "Distribution complete: $DIST_DIR"
