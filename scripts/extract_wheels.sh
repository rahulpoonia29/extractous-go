#!/bin/bash
set -euo pipefail

#=============================================================
# Extractous Binary Downloader & Extractor
#=============================================================
EXTRACTOUS_VERSION="${1:-0.3.0}"
PACKAGE="extractous"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
NATIVE_DIR="$PROJECT_ROOT/native"

# Temporary dirs
DOWNLOAD_DIR=$(mktemp -d -t extractous_download_XXXX)
EXTRACT_DIR=$(mktemp -d -t extractous_extract_XXXX)
trap "rm -rf '$DOWNLOAD_DIR' '$EXTRACT_DIR'" EXIT

echo "==================================================================="
echo "Extractous Binary Fetch Script"
echo "==================================================================="
echo "Version: $EXTRACTOUS_VERSION"
echo "Project root: $PROJECT_ROOT"
echo ""

#-------------------------------------------------------------
# Check dependencies
#-------------------------------------------------------------
for cmd in curl jq unzip; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: Required command '$cmd' not found."
    exit 1
  fi
done

#-------------------------------------------------------------
# Detect platform
#-------------------------------------------------------------
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux*)   PLATFORM="linux"; LIB_EXT="so" ;;
  darwin*)  PLATFORM="macosx"; LIB_EXT="dylib" ;;
  mingw*|msys*|cygwin*|windows*) PLATFORM="win"; LIB_EXT="dll" ;;
  *) echo "ERROR: Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_TAG="x86_64"; ARCH="amd64" ;;
  aarch64|arm64) ARCH_TAG="arm64"; ARCH="arm64" ;;
  *) echo "ERROR: Unsupported architecture: $ARCH"; exit 1 ;;
esac

#-------------------------------------------------------------
# Determine wheel platform tag
#-------------------------------------------------------------
case "${PLATFORM}_${ARCH}" in
  linux_amd64)   WHEEL_TAG="manylinux_2_28_x86_64" ;;
  darwin_amd64)  WHEEL_TAG="macosx_10_12_x86_64" ;;
  darwin_arm64)  WHEEL_TAG="macosx_11_0_arm64" ;;
  win_amd64)     WHEEL_TAG="win_amd64" ;;
  *) echo "ERROR: No wheel available for this platform"; exit 1 ;;
esac

TARGET_DIR="$NATIVE_DIR/${PLATFORM}_${ARCH}"
mkdir -p "$TARGET_DIR"

echo "Detected platform: ${PLATFORM}_${ARCH}"
echo "Wheel tag: $WHEEL_TAG"
echo "Target dir: $TARGET_DIR"
echo ""

#-------------------------------------------------------------
# Fetch wheel URL from PyPI JSON API
#-------------------------------------------------------------
JSON_URL="https://pypi.org/pypi/${PACKAGE}/${EXTRACTOUS_VERSION}/json"
echo "Fetching wheel URL from $JSON_URL ..."

WHEEL_URL=$(curl -s "$JSON_URL" \
  | jq -r --arg TAG "$WHEEL_TAG" '.urls[] | select(.filename | contains($TAG)) | .url' \
  | head -n 1)

if [ -z "$WHEEL_URL" ] || [ "$WHEEL_URL" = "null" ]; then
  echo "ERROR: Could not find wheel for platform tag '$WHEEL_TAG'"
  exit 1
fi

echo "Found wheel: $WHEEL_URL"
echo ""

#-------------------------------------------------------------
# Download the wheel
#-------------------------------------------------------------
WHEEL_FILE="$DOWNLOAD_DIR/$(basename "$WHEEL_URL")"
echo "Downloading wheel..."
curl -L -o "$WHEEL_FILE" "$WHEEL_URL"
echo "Downloaded: $WHEEL_FILE"
echo ""

#-------------------------------------------------------------
# Extract wheel contents
#-------------------------------------------------------------
echo "Extracting wheel..."
unzip -q "$WHEEL_FILE" -d "$EXTRACT_DIR"

echo "Copying libraries..."
LIBS_COPIED=0

for lib in "$EXTRACT_DIR"/extractous/*."${LIB_EXT}"*; do
  if [ -f "$lib" ]; then
    cp "$lib" "$TARGET_DIR/"
    echo "  ✓ $(basename "$lib")"
    LIBS_COPIED=$((LIBS_COPIED + 1))
  fi
done

echo ""
if [ $LIBS_COPIED -eq 0 ]; then
  echo "ERROR: No libraries were extracted! Check wheel contents."
  exit 1
fi

#-------------------------------------------------------------
# Finish
#-------------------------------------------------------------
echo "==================================================================="
echo "✓ SUCCESS: Extracted $LIBS_COPIED libraries"
echo "==================================================================="
echo ""
echo "Libraries installed to: $TARGET_DIR"
ls -lh "$TARGET_DIR"
echo ""
