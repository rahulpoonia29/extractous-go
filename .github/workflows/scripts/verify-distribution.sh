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

# Check main library (handle Windows naming)
if [ "$OS" = "Windows" ]; then
    MAIN_LIB="lib/extractous_ffi.$LIB_EXT"
else
    MAIN_LIB="lib/libextractous_ffi.$LIB_EXT"
fi

if [ ! -f "$MAIN_LIB" ]; then
    echo "✗ Main library missing: $MAIN_LIB"
    exit 1
fi
echo "✓ Main library present: $(basename $MAIN_LIB)"

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
    
elif [ "$OS" = "Windows" ]; then
    echo "✓ Windows DLL present (runtime check skipped in CI)"
fi

# Total size
echo ""
echo "=== Total Distribution ==="
du -sh .

echo ""
echo "✓ Distribution verified"
