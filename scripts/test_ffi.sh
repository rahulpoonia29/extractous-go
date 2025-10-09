#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SMOKE_DIR="$ROOT_DIR/tests/smoke"

echo "==> Running C smoke tests..."

# Build smoke test
cd "$SMOKE_DIR"
make clean
make

# Detect platform
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS-$ARCH" in
    Linux-x86_64)   PLATFORM="linux_amd64" ;;
    Linux-aarch64)  PLATFORM="linux_arm64" ;;
    Darwin-x86_64)  PLATFORM="darwin_amd64" ;;
    Darwin-arm64)   PLATFORM="darwin_arm64" ;;
    *)
        echo "Error: Unsupported platform: $OS-$ARCH"
        exit 1
        ;;
esac

NATIVE_DIR="$ROOT_DIR/native/$PLATFORM"

# Set library path for runtime
if [ "$OS" = "Darwin" ]; then
    export DYLD_LIBRARY_PATH="$NATIVE_DIR:${DYLD_LIBRARY_PATH:-}"
else
    export LD_LIBRARY_PATH="$NATIVE_DIR:${LD_LIBRARY_PATH:-}"
fi

# Copy test files
TEST_DIR="$SMOKE_DIR/test_files"
mkdir -p "$TEST_DIR"
cp "$ROOT_DIR/tests/testdata"/* "$TEST_DIR/" 2>/dev/null || true

# Run smoke test
echo ""
echo "==> Running smoke test binary..."
"$SMOKE_DIR/smoke"

echo ""
echo "==> All tests passed!"
