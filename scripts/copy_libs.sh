#!/bin/bash
set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <platform>"
    echo "Example: $0 linux_amd64"
    exit 1
fi

PLATFORM=$1
SRC_DIR="extractous-ffi/target/release"
DEST_DIR="native/$PLATFORM"

mkdir -p "$DEST_DIR"

# Copy the FFI library
if [ -f "$SRC_DIR/libextractous_ffi.so" ]; then
    cp "$SRC_DIR/libextractous_ffi.so" "$DEST_DIR/"
    echo "Copied libextractous_ffi.so"
elif [ -f "$SRC_DIR/libextractous_ffi.dylib" ]; then
    cp "$SRC_DIR/libextractous_ffi.dylib" "$DEST_DIR/"
    echo "Copied libextractous_ffi.dylib"
elif [ -f "$SRC_DIR/extractous_ffi.dll" ]; then
    cp "$SRC_DIR/extractous_ffi.dll" "$DEST_DIR/"
    echo "Copied extractous_ffi.dll"
fi

# Copy all GraalVM libs (should already be there from extractous-core build)
find "$SRC_DIR" -type f \( -name "*.so" -o -name "*.dylib" -o -name "*.dll" \) \
    ! -name "libextractous_ffi.*" ! -name "extractous_ffi.dll" \
    -exec cp {} "$DEST_DIR/" \;

echo "Libraries copied to $DEST_DIR"
ls -lh "$DEST_DIR"
