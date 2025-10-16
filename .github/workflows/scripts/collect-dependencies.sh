#!/usr/bin/env bash
set -e

# Collect FFI library dependencies recursively
# Usage: ./collect-dependencies.sh <PLATFORM> <TARGET> <LIB_EXT> <OS>

PLATFORM=$1
TARGET=$2
LIB_EXT=$3
OS=$4

echo "=== Dependency build and collection ==="
echo "Platform: $PLATFORM | Target: $TARGET | Ext: $LIB_EXT"
echo ""

mkdir -p "dist/$PLATFORM/lib"
mkdir -p "dist/$PLATFORM/include"

TIKA_NATIVE_DIR=$(find ./ffi/target/"$TARGET"/release/build -type f -path "*extractous*/out/tika-native/build.gradle" | head -1 | xargs dirname)

if [ -z "$TIKA_NATIVE_DIR" ]; then
	echo "✗ Error: Could not find Tika Native Directory in build outputs"
fi

# 1. Build Tika Native with libraries linked
echo "=== Building Tika Native ==="
(
	cd "$TIKA_NATIVE_DIR"
	./gradlew nativeCompile \
		-PnativeImageArgs="--static-nolibc -H:+JNI -H:NativeLinkerOption=-ldl -H:NativeLinkerOption=-lpthread -H:NativeLinkerOption=-lrt"
)

# 2. Copy FFI library
printf "\n=== Copying FFI Library ==="

if [ "$OS" = "Windows" ]; then
	MAIN_LIB="./ffi/target/$TARGET/release/extractous_ffi.$LIB_EXT"
else
	MAIN_LIB="./ffi/target/$TARGET/release/libextractous_ffi.$LIB_EXT"
fi

if [ -f "$MAIN_LIB" ]; then
	cp "$MAIN_LIB" "dist/$PLATFORM/lib/"
	echo "✓ Copied $(basename "$MAIN_LIB")"
else
	echo "✗ Error: Main library not found at $MAIN_LIB"
	exit 1
fi

# 3. Copy libtika_native and its dependencies
printf "\n=== Copying libtika_native and dependencies ==="

LIB_DIR="$TIKA_NATIVE_DIR/build/native/nativeCompile"

if [ -d "$LIB_DIR" ]; then
	# Count libraries first
	LIB_COUNT=$(find "$LIB_DIR" -maxdepth 1 -name "*.$LIB_EXT" -type f 2>/dev/null | wc -l)

	if [ "$LIB_COUNT" -gt 0 ]; then
		# Copy all libraries
		cp "$LIB_DIR"/*."$LIB_EXT" "dist/$PLATFORM/lib/" 2>/dev/null || true
		echo "✓ Copied $LIB_COUNT libraries from $LIB_DIR"

		# Show copied files
		printf "\nCopied files:"
		find "dist/$PLATFORM/lib/"*."$LIB_EXT"
	else
		echo ""
		echo "⚠ Warning: No .$LIB_EXT files found in $LIB_DIR"
	fi
else
	echo "✗ Error: Library directory not found: $LIB_DIR"
	exit 1
fi

# 4. Copy header
echo ""
echo "=== Copy C Header ==="

if [ -f "./ffi/include/extractous.h" ]; then
	cp "./ffi/include/extractous.h" "dist/$PLATFORM/include/"
	echo "✓ Copied extractous.h"
else
	echo "⚠ Warning: extractous.h not found"
fi
