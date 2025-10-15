#!/usr/bin/env bash

set -e

# Collect FFI library dependencies recursively
# Usage: ./collect-dependencies.sh <PLATFORM> <TARGET> <LIB_EXT> <OS>

PLATFORM=$1
TARGET=$2
LIB_EXT=$3
OS=$4

echo "=== Dependency Collection ==="
echo "Platform: $PLATFORM | Target: $TARGET | Ext: $LIB_EXT"
echo ""

# Create distribution structure
mkdir -p "dist/$PLATFORM/lib"
mkdir -p "dist/$PLATFORM/include"

# Find the extractous build directory that contains actual libraries
echo "Searching for libtika_native.$LIB_EXT..."
TIKA_LIB=$(find "./ffi/target/$TARGET/release/build" -type f -name "*tika_native.$LIB_EXT" 2>/dev/null | head -1)

if [ -z "$TIKA_LIB" ]; then
    echo "✗ Error: Could not find libtika_native.$LIB_EXT in build outputs"
    echo "Searching for any extractous build directories..."
    find "./ffi/target/$TARGET/release/build" -maxdepth 2 -name "extractous-*" -type d
    exit 1
fi

LIB_DIR=$(dirname "$TIKA_LIB")
BUILD_DIR=$(dirname "$(dirname "$LIB_DIR")")

echo "Build directory: $BUILD_DIR"
echo "Libraries directory: $LIB_DIR"
echo ""

# 1. Copy main FFI library
echo "=== Main FFI Library ==="

if [ "$OS" = "Windows" ]; then
    MAIN_LIB="./ffi/target/$TARGET/release/extractous_ffi.$LIB_EXT"
else
    MAIN_LIB="./ffi/target/$TARGET/release/libextractous_ffi.$LIB_EXT"
fi

if [ -f "$MAIN_LIB" ]; then
    cp "$MAIN_LIB" "dist/$PLATFORM/lib/"
    echo "✓ Copied $(basename $MAIN_LIB)"
else
    echo "✗ Error: Main library not found at $MAIN_LIB"
    exit 1
fi

# 2. Copy libtika_native and its dependencies
echo ""
echo "=== libtika_native and Dependencies ==="

if [ -d "$LIB_DIR" ]; then
    # Count libraries first
    LIB_COUNT=$(find "$LIB_DIR" -maxdepth 1 -name "*.$LIB_EXT" -type f 2>/dev/null | wc -l)
    
    if [ "$LIB_COUNT" -gt 0 ]; then
        # Copy all libraries
        cp "$LIB_DIR"/*."$LIB_EXT" "dist/$PLATFORM/lib/" 2>/dev/null || true
        echo "✓ Copied $LIB_COUNT libraries from $LIB_DIR"
        
        # Show what we copied
        echo ""
        echo "Copied files:"
        ls -lh "dist/$PLATFORM/lib/"*."$LIB_EXT" 2>/dev/null

        cd "dist/$PLATFORM/lib/"
        for file in *; do if [ -f "$file" ]; then ldd "$file"; fi; done
    else
        echo "⚠ Warning: No .$LIB_EXT files found in $LIB_DIR"
    fi
else
    echo "✗ Error: Library directory not found: $LIB_DIR"
    exit 1
fi

# 3. Copy header
echo ""
echo "=== C Header ==="

if [ -f "./ffi/include/extractous.h" ]; then
    cp "./ffi/include/extractous.h" "dist/$PLATFORM/include/"
    echo "✓ Copied extractous.h"
else
    echo "⚠ Warning: extractous.h not found"
fi

# Summary
echo ""
echo "=== Summary ==="
TOTAL_LIBS=$(find "dist/$PLATFORM/lib" -name "*.$LIB_EXT" -type f 2>/dev/null | wc -l)
TOTAL_SIZE=$(du -sh "dist/$PLATFORM/lib" 2>/dev/null | cut -f1)

echo "Libraries: $TOTAL_LIBS"
echo "Size: $TOTAL_SIZE"
echo ""
echo "✓ Collection complete"