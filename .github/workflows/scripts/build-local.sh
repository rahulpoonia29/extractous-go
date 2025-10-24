#!/usr/bin/env bash
set -e

# Local build script for extractous-ffi
# Usage: ./scripts/build-local.sh [target]

# Detect platform
if [ -z "$1" ]; then
    case "$(uname -s)" in
        Linux*)     TARGET="x86_64-unknown-linux-gnu"; LIB_EXT="so"; PLATFORM="linux_amd64";;
        Darwin*)    
            if [ "$(uname -m)" = "arm64" ]; then
                TARGET="aarch64-apple-darwin"; PLATFORM="darwin_arm64"
            else
                TARGET="x86_64-apple-darwin"; PLATFORM="darwin_amd64"
            fi
            LIB_EXT="dylib"
            ;;
        MINGW*|MSYS*|CYGWIN*) TARGET="x86_64-pc-windows-msvc"; LIB_EXT="dll"; PLATFORM="windows_amd64";;
        *)          echo "Unknown platform"; exit 1;;
    esac
else
    TARGET=$1
    # Detect lib_ext and platform from target
fi

echo "Building for: $TARGET"
echo "Platform: $PLATFORM"

# Check for GraalVM
if [ -z "$GRAALVM_HOME" ] && [ -z "$JAVA_HOME" ]; then
    echo "Error: GRAALVM_HOME or JAVA_HOME must be set"
    echo "Install GraalVM 23+ with native-image"
    exit 1
fi

# Build
cd ffi
cargo build --release --target "$TARGET"
cd ..

# Collect libraries
./scripts/collect-libs.sh "$PLATFORM" "$TARGET" "$LIB_EXT"

echo ""
echo "✓ Build complete!"
echo "Distribution: dist/$PLATFORM/"
echo ""
echo "To use in Go:"
echo "  export CGO_CFLAGS=\"-I$(pwd)/dist/$PLATFORM/include\""
echo "  export CGO_LDFLAGS=\"-L$(pwd)/dist/$PLATFORM/lib -lextractous_ffi\""
echo "  export LD_LIBRARY_PATH=\"$(pwd)/dist/$PLATFORM/lib\" # Linux"
echo "  export DYLD_LIBRARY_PATH=\"$(pwd)/dist/$PLATFORM/lib\" # macOS"
