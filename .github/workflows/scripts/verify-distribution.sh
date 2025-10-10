#!/usr/bin/env bash
set -e

# Verify distribution artifacts
# Usage: ./verify-distribution.sh <platform> <lib_ext>

PLATFORM=$1
LIB_EXT=$2

echo "=== Distribution Verification Script ==="
echo "Platform: $PLATFORM"
echo "Library extension: $LIB_EXT"
echo ""

cd "dist/$PLATFORM"

echo "=== Verifying distribution artifacts ==="

if [ ! -f "lib/libextractous_ffi.$LIB_EXT" ]; then
  echo "✗ Main library missing from distribution!"
  exit 1
fi

echo "✓ Main library present"
echo "Library size: $(du -h lib/libextractous_ffi.$LIB_EXT | cut -f1)"

echo ""
echo "=== All libraries in distribution ==="
find lib -name "*.$LIB_EXT" -type f -exec basename {} \; | sort

echo ""
echo "=== Library details ==="
ls -lh lib/*.$LIB_EXT 2>/dev/null || ls -lh lib/

echo ""
echo "=== Total distribution size ==="
du -sh .

echo ""
echo "✓ Distribution verification completed successfully"
