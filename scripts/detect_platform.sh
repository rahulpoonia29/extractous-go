#!/usr/bin/env bash
# Detect current platform and output in our naming convention

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS-$ARCH" in
    Linux-x86_64)       echo "linux_amd64" ;;
    Darwin-x86_64)      echo "darwin_amd64" ;;
    Darwin-arm64)       echo "darwin_arm64" ;;
    MINGW*-x86_64)      echo "windows_amd64" ;;
    MSYS*-x86_64)       echo "windows_amd64" ;;
    *)
        echo "Error: Unsupported platform: $OS-$ARCH" >&2
        exit 1
        ;;
esac
