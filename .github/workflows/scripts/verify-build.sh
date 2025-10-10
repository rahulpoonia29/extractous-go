#!/usr/bin/env bash
set +e

# Verify FFI build outputs
# Usage: ./verify-build.sh <target> <lib_ext> <os>

TARGET=$1
LIB_EXT=$2
OS=$3

echo "=== Build Verification Script ==="
echo "Target: $TARGET"
echo "Library extension: $LIB_EXT"
echo "OS: $OS"
echo ""

echo "=== Main library verification ==="
MAIN_LIB="./ffi/target/$TARGET/release/libextractous_ffi.$LIB_EXT"

if [ -f "$MAIN_LIB" ]; then
  echo "✓ Found main library: $MAIN_LIB"
  ls -lh "$MAIN_LIB"
  
  # Check dependencies based on platform
  echo ""
  if [ "$OS" = "Linux" ]; then
    echo "Checking dependencies (ldd):"
    ldd "$MAIN_LIB" || true
  elif [ "$OS" = "macOS" ]; then
    echo "Checking dependencies (otool):"
    otool -L "$MAIN_LIB" || true
  elif [ "$OS" = "Windows" ]; then
    echo "Checking dependencies (dumpbin):"
    where dumpbin && dumpbin /dependents "$MAIN_LIB" || echo "dumpbin not available"
  fi
else
  echo "✗ Main library not found: $MAIN_LIB"
  exit 1
fi

echo ""
echo "=== Bundled libraries count ==="
LIB_COUNT=$(find "target/$TARGET/release/build" -name "*.$LIB_EXT" -type f 2>/dev/null | wc -l)
echo "Total bundled libraries: $LIB_COUNT"

echo ""
echo "✓ Build verification completed successfully"
