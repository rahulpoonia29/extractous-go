#!/usr/bin/env bash
set -e

# Collect FFI library dependencies recursively
# Usage: ./collect-dependencies.sh <platform> <target> <lib_ext> <os>

PLATFORM=$1
TARGET=$2
LIB_EXT=$3
OS=$4

echo "=== Dependency Collection Script ==="
echo "Platform: $PLATFORM"
echo "Target: $TARGET"
echo "Library extension: $LIB_EXT"
echo "OS: $OS"
echo ""

mkdir -p "dist/$PLATFORM/lib"

# Find the extractous build directory
BUILD_DIR=$(find "./ffi/target/$TARGET/release/build" -maxdepth 1 -name "extractous-*" -type d | head -1)

if [ -z "$BUILD_DIR" ]; then
  echo "Error: Could not find extractous build directory"
  exit 1
fi

echo "Using build directory: $BUILD_DIR"

# Define paths to search for custom libraries
SEARCH_PATHS=(
  "./ffi/target/$TARGET/release"
  "$BUILD_DIR/out/tika-native/build/native/nativeCompile"
  "$BUILD_DIR/out/libs"
)

# Add GraalVM path if it exists
GRAALVM_PATH=$(find "$BUILD_DIR/out" -type d -name "graalvm-*" 2>/dev/null | head -1)
if [ -n "$GRAALVM_PATH" ]; then
  if [ "$OS" = "Windows" ]; then
    SEARCH_PATHS+=("$GRAALVM_PATH/bin")
  else
    SEARCH_PATHS+=("$GRAALVM_PATH/lib")
  fi
fi

echo "Search paths configured:"
for path in "${SEARCH_PATHS[@]}"; do
  echo "  - $path"
done
echo ""

# Platform-specific system library paths to exclude
if [ "$OS" = "Linux" ]; then
  SYSTEM_PATHS=("/lib/" "/usr/lib/" "/lib64/" "/usr/lib64/" "linux-vdso" "ld-linux")
elif [ "$OS" = "macOS" ]; then
  SYSTEM_PATHS=("/usr/lib/" "/System/" "@rpath" "@loader_path")
elif [ "$OS" = "Windows" ]; then
  SYSTEM_PATHS=("C:\\\\Windows" "C:/Windows" "KERNEL32" "msvcrt")
fi

# Function to check if a library path is a system library
is_system_lib() {
  local lib_path="$1"
  for sys_path in "${SYSTEM_PATHS[@]}"; do
    if [[ "$lib_path" == *"$sys_path"* ]]; then
      return 0
    fi
  done
  return 1
}

# Function to find a library in our search paths
find_custom_lib() {
  local lib_name="$1"
  
  # First, try direct file match in each search path
  for search_path in "${SEARCH_PATHS[@]}"; do
    if [ -f "$search_path/$lib_name" ]; then
      echo "$search_path/$lib_name"
      return 0
    fi
  done
  
  # Then, search recursively (slower but more thorough)
  for search_path in "${SEARCH_PATHS[@]}"; do
    local found=$(find "$search_path" -name "$lib_name" -type f 2>/dev/null | head -1)
    if [ -n "$found" ]; then
      echo "$found"
      return 0
    fi
  done
  
  return 1
}

# Array to track processed libraries
declare -A processed_libs

# Platform-specific dependency checking function
get_dependencies() {
  local lib_path="$1"
  
  if [ "$OS" = "Linux" ]; then
    # Use ldd to get dependencies
    ldd "$lib_path" 2>/dev/null | grep '=>' | awk '{print $1, $3}'
  elif [ "$OS" = "macOS" ]; then
    otool -L "$lib_path" 2>/dev/null | tail -n +2 | awk '{print $1, $1}'
  elif [ "$OS" = "Windows" ]; then
    if command -v dumpbin &> /dev/null; then
      dumpbin /dependents "$lib_path" 2>/dev/null | grep "\.dll" | awk '{print $1, $1}'
    else
      echo ""
    fi
  fi
}

# Function to recursively collect dependencies
collect_deps() {
  local lib_path="$1"
  local lib_name=$(basename "$lib_path")
  
  # Skip if already processed
  if [[ -n "${processed_libs[$lib_name]}" ]]; then
    echo "  (Already processed: $lib_name)"
    return
  fi
  
  echo "Processing: $lib_name"
  processed_libs[$lib_name]=1
  
  # Copy the library
  cp "$lib_path" "dist/$PLATFORM/lib/"
  echo "  ✓ Copied: $lib_name"
  
  # Get dependencies
  local deps=$(get_dependencies "$lib_path")
  
  if [ -z "$deps" ]; then
    echo "  - No dependencies found or all are system libraries"
    return
  fi
  
  # Process each dependency
  while IFS= read -r line; do
    if [ -z "$line" ]; then
      continue
    fi
    
    local dep_name=$(echo "$line" | awk '{print $1}')
    local dep_resolved_path=$(echo "$line" | awk '{print $2}')
    
    # Skip if no dependency name
    if [ -z "$dep_name" ] || [ "$dep_name" = "not" ]; then
      continue
    fi
    
    echo "  → Checking dependency: $dep_name"
    
    # Skip if it's a system library (check resolved path)
    if [ -n "$dep_resolved_path" ] && is_system_lib "$dep_resolved_path"; then
      echo "    ✗ Skipping system lib: $dep_name (at $dep_resolved_path)"
      continue
    fi
    
    # Check if dependency is already in our distribution
    if [ -f "dist/$PLATFORM/lib/$dep_name" ]; then
      echo "    ✓ Already in distribution: $dep_name"
      continue
    fi
    
    # Try to find the library in our build directories
    local custom_lib=$(find_custom_lib "$dep_name")
    
    if [ -n "$custom_lib" ]; then
      echo "    ✓ Found in build: $custom_lib"
      # Recursively process this dependency
      collect_deps "$custom_lib"
    else
      # If we couldn't find it in build dirs, check if ldd gave us a resolved path
      if [ -n "$dep_resolved_path" ] && [ -f "$dep_resolved_path" ] && ! is_system_lib "$dep_resolved_path"; then
        echo "    ! Using resolved path: $dep_resolved_path"
        collect_deps "$dep_resolved_path"
      else
        echo "    ! Warning: Dependency $dep_name not found in build directories"
        echo "              Resolved path: ${dep_resolved_path:-not resolved}"
      fi
    fi
  done <<< "$deps"
}

# Start from the main library
MAIN_LIB="./ffi/target/$TARGET/release/libextractous_ffi.$LIB_EXT"

if [ ! -f "$MAIN_LIB" ]; then
  echo "Error: Main library not found: $MAIN_LIB"
  exit 1
fi

echo "=== Starting recursive dependency collection ==="
echo "Main library: $MAIN_LIB"
echo ""

collect_deps "$MAIN_LIB"

echo ""
echo "=== Collection complete ==="
echo "Libraries collected:"
ls -1 "dist/$PLATFORM/lib/" | sort

echo ""
echo "=== Library details ==="
ls -lh "dist/$PLATFORM/lib/"

echo ""
echo "=== Dependency verification ==="
for lib in dist/$PLATFORM/lib/*.$LIB_EXT; do
  lib_name=$(basename "$lib")
  echo "Checking $lib_name:"
  if [ "$OS" = "Linux" ]; then
    ldd "$lib" 2>/dev/null | grep "not found" && echo "  ⚠ Missing dependencies!" || echo "  ✓ All dependencies resolved"
  elif [ "$OS" = "macOS" ]; then
    otool -L "$lib" 2>/dev/null | tail -n +2 | while read -r line; do
      dep=$(echo "$line" | awk '{print $1}')
      if [[ "$dep" != /usr/* ]] && [[ "$dep" != /System/* ]] && [[ "$dep" != @* ]]; then
        if [ ! -f "dist/$PLATFORM/lib/$(basename "$dep")" ]; then
          echo "  ⚠ Missing: $(basename "$dep")"
        fi
      fi
    done
  fi
done

echo ""
echo "=== Total size ==="
du -sh "dist/$PLATFORM/lib/"

# Copy headers if they exist
if [ -d "include" ]; then
  echo ""
  echo "Copying header files..."
  cp -r include "dist/$PLATFORM/"
fi

echo ""
echo "✓ Dependency collection completed successfully"
