#!/bin/bash
set -e

EXTRACTOUS_VERSION="0.3.0"
BASE_URL="https://files.pythonhosted.org/packages"

declare -A WHEELS=(
    ["linux_amd64"]="6f/d9/1a3838e24f78902ca1a594110a00812134102c5cad13f889141509062481/extractous-${EXTRACTOUS_VERSION}-cp38-abi3-manylinux_2_28_x86_64.whl"
    ["darwin_arm64"]="66/91/7debbfabadb88d34687bf93e23d176692bdae7e82c51180b2481710bb709/extractous-${EXTRACTOUS_VERSION}-cp38-abi3-macosx_11_0_arm64.whl"
    ["darwin_amd64"]="98/50/99d6e8982ced454cc7a0e184988b63c65e199587626c45404fc7b6ab9d90/extractous-${EXTRACTOUS_VERSION}-cp38-abi3-macosx_10_12_x86_64.whl"
    ["windows_amd64"]="07/a1/dd01a3abb4c4af89cf3775735948d76522233ae3550a166b8c2f7c849a52/extractous-${EXTRACTOUS_VERSION}-cp38-abi3-win_amd64.whl"
)

echo "Extracting native libraries from Python wheels v${EXTRACTOUS_VERSION}..."

for platform in "${!WHEELS[@]}"; do
    echo "  Processing $platform..."
    
    url="${BASE_URL}/${WHEELS[$platform]}"
    wheel_name=$(basename "$url")
    
    wget -q "$url" -O "/tmp/$wheel_name"
    unzip -q "/tmp/$wheel_name" -d "/tmp/$platform"
    
    mkdir -p "native/$platform"
    
    if [[ "$platform" == windows* ]]; then
        cp /tmp/$platform/extractous/*.dll native/$platform/
        rm -f native/$platform/_extractous*.pyd
    else
        find /tmp/$platform/extractous -type f \( -name "*.so" -o -name "*.dylib" \) -exec cp {} native/$platform/ \;
        rm -f native/$platform/_extractous*.so
    fi
    
    count=$(ls native/$platform | wc -l)
    echo "    ✓ Extracted $count libraries"
    
    rm -rf "/tmp/$platform" "/tmp/$wheel_name"
done

echo "All native libraries extracted successfully!"
