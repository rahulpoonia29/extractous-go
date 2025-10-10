#!/bin/bash

# Build native libraries for all supported platforms
# This script should be run before creating a new release

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/releases"

echo "Building extractous-ffi for all platforms..."
echo "Output directory: $OUTPUT_DIR"

# Create output directory
mkdir -p "$OUTPUT_DIR"

cd "$PROJECT_ROOT/ffi"

# Array of target platforms
declare -a TARGETS=(
    "x86_64-unknown-linux-gnu:linux_amd64:tar.gz"
    "aarch64-unknown-linux-gnu:linux_arm64:tar.gz"
    "x86_64-apple-darwin:darwin_amd64:tar.gz"
    "aarch64-apple-darwin:darwin_arm64:tar.gz"
    "x86_64-pc-windows-gnu:windows_amd64:zip"
)

# Build for each target
for target_spec in "${TARGETS[@]}"; do
    IFS=':' read -r rust_target platform_dir archive_type <<< "$target_spec"
    
    echo ""
    echo "========================================="
    echo "Building for: $platform_dir"
    echo "Rust target: $rust_target"
    echo "========================================="
    
    # Build
    echo "Running: cargo build --release --target $rust_target"
    cargo build --release --target "$rust_target"
    
    # Create staging directory
    STAGE_DIR="$OUTPUT_DIR/stage/$platform_dir"
    mkdir -p "$STAGE_DIR"
    
    # Copy libraries based on platform
    echo "Copying libraries to $STAGE_DIR..."
    
    case "$platform_dir" in
        windows_*)
            # Windows: .dll files
            cp "target/$rust_target/release/extractous_ffi.dll" "$STAGE_DIR/" 2>/dev/null || \
            cp "target/$rust_target/release/libextractous_ffi.dll" "$STAGE_DIR/extractous_ffi.dll" 2>/dev/null || true
            
            # Copy any additional Windows DLLs
            find "../native/$platform_dir" -name "*.dll" -exec cp {} "$STAGE_DIR/" \; 2>/dev/null || true
            ;;
        darwin_*)
            # macOS: .dylib files
            cp "target/$rust_target/release/libextractous_ffi.dylib" "$STAGE_DIR/" 2>/dev/null || true
            
            # Copy any additional macOS dylibs
            find "../native/$platform_dir" -name "*.dylib" -exec cp {} "$STAGE_DIR/" \; 2>/dev/null || true
            ;;
        linux_*)
            # Linux: .so files
            cp "target/$rust_target/release/libextractous_ffi.so" "$STAGE_DIR/" 2>/dev/null || true
            
            # Copy any additional Linux .so files
            find "../native/$platform_dir" -name "*.so" -exec cp {} "$STAGE_DIR/" \; 2>/dev/null || true
            ;;
    esac
    
    # Create archive
    ARCHIVE_NAME="extractous-ffi-$platform_dir"
    
    echo "Creating archive: $ARCHIVE_NAME.$archive_type"
    
    cd "$STAGE_DIR"
    
    if [ "$archive_type" = "tar.gz" ]; then
        tar -czf "$OUTPUT_DIR/$ARCHIVE_NAME.tar.gz" ./*
    elif [ "$archive_type" = "zip" ]; then
        zip -r "$OUTPUT_DIR/$ARCHIVE_NAME.zip" ./*
    fi
    
    cd "$PROJECT_ROOT/ffi"
    
    echo "✓ Created: $OUTPUT_DIR/$ARCHIVE_NAME.$archive_type"
    
    # Calculate size
    ARCHIVE_PATH="$OUTPUT_DIR/$ARCHIVE_NAME.$archive_type"
    if [ -f "$ARCHIVE_PATH" ]; then
        SIZE=$(du -h "$ARCHIVE_PATH" | cut -f1)
        echo "  Size: $SIZE"
    fi
done

# Cleanup staging
rm -rf "$OUTPUT_DIR/stage"

echo ""
echo "========================================="
echo "Build complete!"
echo "========================================="
echo "Archives created in: $OUTPUT_DIR"
echo ""
echo "Next steps:"
echo "1. Create a new GitHub release (e.g., v0.2.1)"
echo "2. Upload all .tar.gz and .zip files from $OUTPUT_DIR"
echo "3. Update LibraryVersion in internal/download/download.go"
echo ""
ls -lh "$OUTPUT_DIR"
