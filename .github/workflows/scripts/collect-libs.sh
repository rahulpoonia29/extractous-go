#!/usr/bin/env bash
set -e

# Collect all native libraries for distribution
# Usage: ./scripts/collect-libs.sh <platform> <target> <lib_ext>

PLATFORM=$1 # e.g., linux_amd64
TARGET=$2   # e.g., x86_64-unknown-linux-gnu
LIB_EXT=$3  # e.g., so, dll, dylib

echo "=== Collecting Libraries ==="
echo "Platform: $PLATFORM | Target: $TARGET | Extension: $LIB_EXT"

# Create distribution structure
DIST_DIR="dist/$PLATFORM"
mkdir -p "$DIST_DIR"
# Determine OS-specific naming
if [[ "$PLATFORM" == *"windows"* ]]; then
	OS="Windows"
	MAIN_LIB_PREFIX=""
	TIKA_LIB_PREFIX=""
elif [[ "$PLATFORM" == *"darwin"* ]]; then
	OS="macOS"
	MAIN_LIB_PREFIX="lib"
	TIKA_LIB_PREFIX="lib"
else
	OS="Linux"
	MAIN_LIB_PREFIX="lib"
	TIKA_LIB_PREFIX="lib"
fi

# 1. Find and copy main FFI library
RELEASE_DIR="ffi/target/$TARGET/release"
MAIN_LIB="$RELEASE_DIR/${MAIN_LIB_PREFIX}extractous_ffi.$LIB_EXT"

if [ ! -f "$MAIN_LIB" ]; then
	echo "✗ Error: Main FFI library not found: $MAIN_LIB"
	exit 1
fi

echo "✓ Found main FFI library: $MAIN_LIB"
cp "$MAIN_LIB" "$DIST_DIR"

# 2. Find extractous build output directory
# Look for the canonical libs directory created by extractous build.rs
echo ""
echo "Searching for extractous dependencies..."

BUILD_BASE="$RELEASE_DIR/build"

# Find all extractous-*/out/libs directories, sort by modification time (newest first)
LIBS_DIR=$(find "$BUILD_BASE" -type d -path "*/extractous-*/out/libs" -printf "%T@ %p\n" 2>/dev/null |
	sort -rn |
	head -1 |
	cut -d' ' -f2)

# Fallback for macOS (no -printf support)
if [ -z "$LIBS_DIR" ]; then
	# Find all matching directories
	FOUND_DIRS=$(find "$BUILD_BASE" -type d -path "*/extractous-*/out/libs" 2>/dev/null)

	if [ -n "$FOUND_DIRS" ]; then
		# Get the most recently modified directory
		LIBS_DIR=$(echo "$FOUND_DIRS" | while read -r dir; do
			echo "$(stat -f "%m" "$dir") $dir"
		done | sort -rn | head -1 | cut -d' ' -f2-)
	fi
fi

if [ -z "$LIBS_DIR" ]; then
	echo "✗ Error: Could not find extractous out/libs directory"
	echo "Searched in: $BUILD_BASE/extractous-*/out/libs"
	echo ""
	echo "Available build directories:"
	find "$BUILD_BASE" -type d -name "extractous-*" 2>/dev/null || echo "None found"
	echo ""
	echo "Checking for out directories:"
	find "$BUILD_BASE" -type d -name "out" 2>/dev/null || echo "None found"
	exit 1
fi

echo "✓ Found libs directory: $LIBS_DIR"

# 3. Verify libtika_native exists
# Try both with and without prefix for Windows compatibility
TIKA_LIB="$LIBS_DIR/${TIKA_LIB_PREFIX}tika_native.$LIB_EXT"
if [ ! -f "$TIKA_LIB" ] && [ "$OS" = "Windows" ]; then
	# Try with lib prefix on Windows as fallback
	TIKA_LIB_ALT="$LIBS_DIR/libtika_native.$LIB_EXT"
	if [ -f "$TIKA_LIB_ALT" ]; then
		TIKA_LIB="$TIKA_LIB_ALT"
	fi
fi

if [ ! -f "$TIKA_LIB" ]; then
	echo "✗ Error: tika_native.$LIB_EXT not found in $LIBS_DIR"
	echo "Directory contents:"
	ls -lh "$LIBS_DIR" || echo "Directory not accessible"

	# Show all DLL/SO/DYLIB files to help debug
	echo ""
	echo "All native libraries found:"
	find "$LIBS_DIR" -name "*.$LIB_EXT" -o -name "*.dll" -o -name "*.so" -o -name "*.dylib" 2>/dev/null || echo "None found"
	exit 1
fi

echo "✓ Found libtika_native: $TIKA_LIB"

# 4. Copy ALL libraries from out/libs/
# These are all required dependencies bundled by GraalVM
echo ""
echo "Copying all native dependencies..."
cp "$LIBS_DIR"/*."$LIB_EXT" "$DIST_DIR" || {
	echo "✗ Error: Failed to copy libraries"
	exit 1
}

# 5. Patch libextractous_ffi on macOS to use @rpath and replace absolute path
# https://github.com/rahulpoonia29/extractous-go/issues/5
if [ "$OS" = "macOS" ]; then
	echo ""
	echo "Verify XCode tools"    
	# XCode tools are present on github macOS runners by default, but verify anyway
	which otool || { echo "✗ otool not found"; exit 1; }
	which install_name_tool || { echo "✗ install_name_tool not found"; exit 1; }
	otool -L "$DIST_DIR/libextractous_ffi.dylib" || { echo "✗ otool test failed"; exit 1; }

	echo "Patching libextractous_ffi.dylib to use @loader_path for tika"
	OLD_PATH=$(otool -L "$DIST_DIR/libextractous_ffi.dylib" | grep libtika_native.dylib | awk '{print $1}')
	echo "  Old tika_native path: $OLD_PATH"

	# Replace with @loader_path (directory of the main library)
	install_name_tool -change "$OLD_PATH" "@loader_path/libtika_native.dylib" "$DIST_DIR/libextractous_ffi.dylib"

	echo "Debug: New Path"
	echo "!! Should be @loader_path/libtika_native.dylib"
	otool -L "$DIST_DIR/libextractous_ffi.dylib" | sed 's/^/    /'
fi

# Count copied libraries
LIB_COUNT=$(find "$DIST_DIR" -name "*.$LIB_EXT" | wc -l)
echo "✓ Copied $LIB_COUNT libraries"

# 5. Display distribution contents
echo ""
echo "=== Distribution Contents ==="
echo "Libraries ($LIB_COUNT total):"
ls -lh "$DIST_DIR" | tail -n +2

if [ -d "$DIST_DIR/include" ] && [ -n "$(ls -A "$DIST_DIR/include" 2>/dev/null)" ]; then
	echo ""
	echo "Headers:"
	ls -lh "$DIST_DIR/include/"
fi

# 7. Verify dependencies (platform-specific)
echo ""
echo "=== Dependency Verification ==="

if [ "$OS" = "Linux" ]; then
	echo "Checking RPATH configuration..."
	for lib in "$DIST_DIR"*.$LIB_EXT; do
		LIB_NAME=$(basename "$lib")
		echo "  $LIB_NAME:"
		readelf -d "$lib" | grep -E "RPATH|RUNPATH" | sed 's/^/    /' || echo "    No RPATH set"
	done

	echo ""
	echo "Checking main FFI dependencies..."
	ldd "$DIST_DIR/${MAIN_LIB_PREFIX}extractous_ffi.$LIB_EXT" || true

elif [ "$OS" = "macOS" ]; then
	echo "Checking install names..."
	for lib in "$DIST_DIR"*.$LIB_EXT; do
		LIB_NAME=$(basename "$lib")
		echo "  $LIB_NAME:"
		otool -L "$lib" | grep -v "$LIB_NAME:" | sed 's/^/    /'
	done

elif [ "$OS" = "Windows" ]; then
	echo "Windows DLL validation..."
	file "$DIST_DIR"/*.dll 2>/dev/null || echo "  DLLs present"
fi

# 8. Calculate total size
echo ""
echo "=== Distribution Size ==="
echo "Total library size: $(du -sh "$DIST_DIR" | cut -f1)"
echo "Distribution complete: $DIST_DIR"
