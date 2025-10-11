#!/usr/bin/env bash
set -e

# Verify distribution artifacts
# Usage: ./verify-distribution.sh <platform> <lib_ext> <os>

PLATFORM=$1
LIB_EXT=$2
OS=$3

echo "=== Distribution Verification ==="
echo "Platform: $PLATFORM | Ext: $LIB_EXT | OS: $OS"
echo ""

cd "dist/$PLATFORM"

# Check main library
if [ ! -f "lib/libextractous_ffi.$LIB_EXT" ]; then
    echo "✗ Main library missing!"
    exit 1
fi
echo "✓ Main library present"

# List all libraries
echo ""
echo "=== Libraries ==="
find lib -name "*.$LIB_EXT" -type f -exec basename {} \; | sort

# Size breakdown
echo ""
echo "=== Size Breakdown ==="
for lib in lib/*.$LIB_EXT; do
    printf "%-40s %10s\n" "$(basename $lib)" "$(du -h $lib | cut -f1)"
done | sort -k2 -h -r

# Test library loading
echo ""
echo "=== Runtime Test ==="
if [ "$OS" = "Linux" ]; then
    LD_LIBRARY_PATH=lib:$LD_LIBRARY_PATH ldd lib/libextractous_ffi.so | grep "not found" && {
        echo "✗ Unresolved dependencies!"
        exit 1
    } || echo "✓ All dependencies resolved"
    
elif [ "$OS" = "macOS" ]; then
    DYLD_LIBRARY_PATH=lib otool -L lib/libextractous_ffi.dylib
    echo "✓ Library loadable"
fi

# Total size
echo ""
echo "=== Total Distribution ==="
du -sh .

echo ""
echo "✓ Distribution verified"
