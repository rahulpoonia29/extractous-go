#!/bin/bash
set -euo pipefail

#=============================================================
# Extractous - Multi-Platform Binary Extractor
#=============================================================
EXTRACTOUS_VERSION="${1:-0.3.0}"
PACKAGE="extractous"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
NATIVE_DIR="$PROJECT_ROOT/native"

# Temporary directories
DOWNLOAD_DIR=$(mktemp -d -t extractous_download_XXXX)
EXTRACT_DIR_BASE=$(mktemp -d -t extractous_extract_XXXX)
trap "rm -rf '$DOWNLOAD_DIR' '$EXTRACT_DIR_BASE'" EXIT

echo "==================================================================="
echo "Extractous Multi-Platform Extraction Script"
echo "==================================================================="
echo "Version: $EXTRACTOUS_VERSION"
echo "Project root: $PROJECT_ROOT"
echo ""

#-------------------------------------------------------------
# Dependencies check
#-------------------------------------------------------------
for cmd in curl jq unzip sha256sum; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: Required command '$cmd' not found."
    exit 1
  fi
done

#-------------------------------------------------------------
# Fetch package JSON metadata
#-------------------------------------------------------------
JSON_URL="https://pypi.org/pypi/${PACKAGE}/${EXTRACTOUS_VERSION}/json"
echo "Fetching PyPI metadata from $JSON_URL..."
JSON_DATA=$(curl -s "$JSON_URL")

if [ -z "$JSON_DATA" ]; then
  echo "ERROR: Failed to fetch package metadata."
  exit 1
fi

#-------------------------------------------------------------
# Define known platform mappings
#-------------------------------------------------------------
declare -A PLATFORM_MAP=(
  ["manylinux_2_28_x86_64"]="linux_amd64"
  ["macosx_10_12_x86_64"]="darwin_amd64"
  ["macosx_11_0_arm64"]="darwin_arm64"
  ["win_amd64"]="windows_amd64"
)

#-------------------------------------------------------------
# Loop through each wheel entry in JSON
#-------------------------------------------------------------
echo ""
echo "Scanning available wheels..."
echo ""

echo "$JSON_DATA" | jq -c '.urls[] | select(.packagetype=="bdist_wheel")' | while read -r entry; do
  FILENAME=$(echo "$entry" | jq -r '.filename')
  URL=$(echo "$entry" | jq -r '.url')
  SHA256=$(echo "$entry" | jq -r '.digests.sha256')

  # Match the platform
  PLATFORM_TAG=""
  for key in "${!PLATFORM_MAP[@]}"; do
    if [[ "$FILENAME" == *"$key"* ]]; then
      PLATFORM_TAG="$key"
      break
    fi
  done

  if [ -z "$PLATFORM_TAG" ]; then
    echo "Skipping unknown platform wheel: $FILENAME"
    continue
  fi

  PLATFORM_DIR="${PLATFORM_MAP[$PLATFORM_TAG]}"
  TARGET_DIR="$NATIVE_DIR/${PLATFORM_DIR}"
  mkdir -p "$TARGET_DIR"

  echo "-------------------------------------------------------------"
  echo "Processing wheel: $FILENAME"
  echo "Platform: $PLATFORM_DIR"
  echo "URL: $URL"
  echo "Target: $TARGET_DIR"
  echo ""

  # Download
  WHEEL_FILE="$DOWNLOAD_DIR/$FILENAME"
  echo "Downloading..."
  curl -L -s -o "$WHEEL_FILE" "$URL"

  # Verify SHA256 checksum
  echo "$SHA256  $WHEEL_FILE" | sha256sum -c --status || {
    echo "ERROR: SHA256 checksum failed for $FILENAME"
    exit 1
  }
  echo "✓ Checksum verified"

  # Extract
  EXTRACT_DIR="$EXTRACT_DIR_BASE/$PLATFORM_DIR"
  mkdir -p "$EXTRACT_DIR"
  unzip -q "$WHEEL_FILE" -d "$EXTRACT_DIR"

  # Copy libs
  case "$PLATFORM_DIR" in
    linux_amd64)   EXT="so" ;;
    darwin_amd64|darwin_arm64) EXT="dylib" ;;
    windows_amd64) EXT="dll" ;;
  esac

  LIBS_FOUND=0
  for lib in "$EXTRACT_DIR"/extractous/*."$EXT"*; do
    if [ -f "$lib" ]; then
      cp "$lib" "$TARGET_DIR/"
      echo "  ✓ Copied $(basename "$lib")"
      LIBS_FOUND=$((LIBS_FOUND + 1))
    fi
  done

  if [ "$LIBS_FOUND" -eq 0 ]; then
    echo "⚠️  No libraries found in $FILENAME"
  else
    echo "✓ Extracted $LIBS_FOUND files for $PLATFORM_DIR"
  fi

  echo ""
done

#-------------------------------------------------------------
# Summary
#-------------------------------------------------------------
echo "==================================================================="
echo "✅ Extraction completed for all available platforms"
echo "==================================================================="
echo ""
echo "Resulting structure:"
find "$NATIVE_DIR" -type f | sed "s|$PROJECT_ROOT/||"
echo ""
