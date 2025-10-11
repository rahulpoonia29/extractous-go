#!/usr/bin/env bash
set -e

# Collect FFI library dependencies recursively
# Usage: ./collect-dependencies.sh <platform> <target> <lib_ext> <os>

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

# Find extractous build directory
BUILD_DIR=$(find "./ffi/target/$TARGET/release/build" -maxdepth 1 -name "extractous-*" -type d | head -1)

if [ -z "$BUILD_DIR" ]; then
    echo "✗ Error: Could not find extractous build directory"
    exit 1
fi

echo "Build directory: $BUILD_DIR"

# Find GraalVM directory
GRAALVM_DIR="$BUILD_DIR/out"
echo "GraalVM directory: ${GRAALVM_DIR:-Not found}"
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
    echo "✓ Copied libextractous_ffi.$LIB_EXT"
else
    echo "✗ Error: Main library not found"
    exit 1
fi

# 2. Copy libtika_native and its dependencies
echo ""
echo "=== libtika_native and Dependencies ==="
LIB_DIR="$GRAALVM_DIR/tika-native/build/native/nativeCompile"

if [ -f "$LIB_DIR" ]; then
    cp "$LIB_DIR"/*."$LIB_EXT" "dist/$PLATFORM/lib/"
    echo "✓ Copied dependencies from $LIB_DIR"
else
    echo "⚠ Warning: No dependencies found in $LIB_DIR"
fi

# 3. Copy header
echo ""
echo "=== C Header ==="
if [ -f "./ffi/include/extractous.h" ]; then
    cp "./ffi/include/extractous.h" "dist/$PLATFORM/include/"
    echo "✓ Copied extractous.h"
fi

# Summary
echo ""
echo "=== Summary ==="
echo "Libraries: $(find dist/$PLATFORM/lib -name "*.$LIB_EXT" | wc -l)"
echo "Size: $(du -sh dist/$PLATFORM/lib | cut -f1)"
echo ""
echo "✓ Collection complete"
