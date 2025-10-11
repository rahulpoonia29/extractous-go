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
GRAALVM_DIR=$(find "$BUILD_DIR/out" -type d -name "graalvm-*" 2>/dev/null | head -1)
echo "GraalVM directory: ${GRAALVM_DIR:-Not found}"
echo ""

# 1. Copy main FFI library
echo "=== Main FFI Library ==="
MAIN_LIB="./ffi/target/$TARGET/release/libextractous_ffi.$LIB_EXT"

if [ -f "$MAIN_LIB" ]; then
    cp "$MAIN_LIB" "dist/$PLATFORM/lib/"
    echo "✓ Copied libextractous_ffi.$LIB_EXT"
else
    echo "✗ Error: Main library not found"
    exit 1
fi

# 2. Copy libtika_native
echo ""
echo "=== Tika Native Library ==="
TIKA_LIB=$(find "$BUILD_DIR/out/tika-native/build/native/nativeCompile" -name "libtika_native.$LIB_EXT" 2>/dev/null | head -1)

if [ -f "$TIKA_LIB" ]; then
    cp "$TIKA_LIB" "dist/$PLATFORM/lib/"
    echo "✓ Copied libtika_native.$LIB_EXT"
else
    echo "⚠ Warning: libtika_native not found"
fi

# 3. Copy ALL GraalVM runtime libraries directly
echo ""
echo "=== GraalVM Runtime Libraries ==="

if [ -n "$GRAALVM_DIR" ]; then
    if [ "$OS" = "Windows" ]; then
        # Windows: Copy from bin/
        if [ -d "$GRAALVM_DIR/bin" ]; then
            cp "$GRAALVM_DIR/bin"/*.$LIB_EXT "dist/$PLATFORM/lib/" 2>/dev/null || true
            COUNT=$(find "dist/$PLATFORM/lib" -name "*.$LIB_EXT" | wc -l)
            echo "✓ Copied $COUNT DLLs from GraalVM"
        fi
    else
        # Linux/macOS: Copy from lib/
        if [ -d "$GRAALVM_DIR/lib" ]; then
            cp "$GRAALVM_DIR/lib"/*.$LIB_EXT "dist/$PLATFORM/lib/" 2>/dev/null || true
            COUNT=$(find "dist/$PLATFORM/lib" -name "*.$LIB_EXT" | wc -l)
            echo "✓ Copied $COUNT libraries from GraalVM"
        fi
    fi
fi

# 4. Copy header
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
