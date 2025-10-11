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

# Find all extractous build directories and select the most recently modified one
echo "Searching for extractous build directories..."

# Use ls -t for cross-platform sorting by modification time
BUILD_DIR=$(ls -dt ./ffi/target/$TARGET/release/build/extractous-*/ 2>/dev/null | head -1 | sed 's:/$::')

if [ -z "$BUILD_DIR" ]; then
    echo "✗ Error: Could not find any extractous-* build directories"
    echo "Available directories:"
    ls -la "./ffi/target/$TARGET/release/build/" 2>/dev/null || true
    exit 1
fi

echo "Latest build directory: $BUILD_DIR"
echo "Modified: $(stat -c '%y' "$BUILD_DIR" 2>/dev/null || stat -f '%Sm' "$BUILD_DIR" 2>/dev/null || echo "unknown")"

# Find GraalVM libraries directory
GRAALVM_DIR="$BUILD_DIR/out"
LIB_DIR="$GRAALVM_DIR/libs"

echo "GraalVM directory: $GRAALVM_DIR"
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
    # Check if directory has any libraries
    LIB_COUNT=$(ls -1 "$LIB_DIR"/*."$LIB_EXT" 2>/dev/null | wc -l)
    
    if [ "$LIB_COUNT" -gt 0 ]; then
        echo "Found libraries:"
        ls -lh "$LIB_DIR"/*."$LIB_EXT"
        echo ""
        
        # Copy all libraries
        cp "$LIB_DIR"/*."$LIB_EXT" "dist/$PLATFORM/lib/"
        echo "✓ Copied $LIB_COUNT libraries from $LIB_DIR"
    else
        echo "⚠ Warning: No .$LIB_EXT files found in $LIB_DIR"
        echo "Directory contents:"
        ls -la "$LIB_DIR" 2>/dev/null || echo "Directory is empty or inaccessible"
    fi
else
    echo "⚠ Warning: Libraries directory not found: $LIB_DIR"
    echo "Checking build directory structure:"
    find "$BUILD_DIR" -maxdepth 3 -type d 2>/dev/null | head -20
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
echo "Libraries: $(find dist/$PLATFORM/lib -name "*.$LIB_EXT" 2>/dev/null | wc -l)"
echo "Size: $(du -sh dist/$PLATFORM/lib 2>/dev/null | cut -f1)"
echo ""
echo "Final distribution contents:"
ls -lh "
