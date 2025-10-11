#!/usr/bin/env bash
set +e

# Verify FFI build outputs
# Usage: ./verify-build.sh <target> <lib_ext> <os>

TARGET=$1
LIB_EXT=$2
OS=$3

echo "=== Build Verification ==="
echo "Target: $TARGET | Ext: $LIB_EXT | OS: $OS"
echo ""

MAIN_LIB="./ffi/target/$TARGET/release/libextractous_ffi.$LIB_EXT"

if [ -f "$MAIN_LIB" ]; then
    echo "✓ Main library found"
    ls -lh "$MAIN_LIB"
    echo ""
    
    # Check dependencies
    if [ "$OS" = "Linux" ]; then
        echo "Dependencies (ldd):"
        ldd "$MAIN_LIB" || true
        echo ""
        echo "RPATH check:"
        readelf -d "$MAIN_LIB" | grep -E "RPATH|RUNPATH" || echo "No RPATH set"
        
    elif [ "$OS" = "macOS" ]; then
        echo "Dependencies (otool):"
        otool -L "$MAIN_LIB" || true
        echo ""
        echo "RPATH check:"
        otool -l "$MAIN_LIB" | grep -A2 RPATH || echo "No RPATH set"
        
    elif [ "$OS" = "Windows" ]; then
        echo "Library exists (Windows DLL)"
        where dumpbin && dumpbin /dependents "$MAIN_LIB" || echo "dumpbin not available"
    fi
else
    echo "✗ Main library not found: $MAIN_LIB"
    exit 1
fi

echo ""
echo "✓ Build verification complete"
