#!/usr/bin/env bash
set -euo pipefail

#=============================================================================
# Extractous Native Library Extractor
#=============================================================================
# Downloads extractous Python wheels and extracts native libraries
# for cross-platform Go bindings.
#=============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
NATIVE_DIR="$PROJECT_ROOT/native"

# Configuration
EXTRACTOUS_VERSION="${EXTRACTOUS_VERSION:-0.3.0}"
PACKAGE="extractous"
PYPI_JSON_URL="https://pypi.org/pypi/${PACKAGE}/${EXTRACTOUS_VERSION}/json"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Supported platforms
declare -A PLATFORM_MAP=(
    ["manylinux_2_28_x86_64"]="linux_amd64"
    ["macosx_10_12_x86_64"]="darwin_amd64"
    ["macosx_11_0_arm64"]="darwin_arm64"
    ["win_amd64"]="windows_amd64"
)

declare -A PLATFORM_EXTENSIONS=(
    ["linux_amd64"]="so"
    ["darwin_amd64"]="dylib"
    ["darwin_arm64"]="dylib"
    ["windows_amd64"]="dll"
)

# Temporary directories
TEMP_BASE="${TMPDIR:-/tmp}/extractous-$$"
DOWNLOAD_DIR="$TEMP_BASE/downloads"
EXTRACT_DIR="$TEMP_BASE/extract"
mkdir -p "$DOWNLOAD_DIR" "$EXTRACT_DIR"

# Cleanup on exit
cleanup() {
    local exit_code=$?
    if [ $exit_code -ne 0 ]; then
        echo -e "\n${RED}✗ Extraction failed with exit code $exit_code${NC}" >&2
    fi
    rm -rf "$TEMP_BASE"
    exit $exit_code
}
trap cleanup EXIT INT TERM

#=============================================================================
# Helper Functions
#=============================================================================

print_header() {
    echo ""
    echo -e "${CYAN}===================================================================${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}===================================================================${NC}"
}

print_section() {
    echo ""
    echo -e "${BLUE}-------------------------------------------------------------------${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}-------------------------------------------------------------------${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1" >&2
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_info() {
    echo -e "${CYAN}→${NC} $1"
}

check_dependencies() {
    local missing=()
    for cmd in curl jq unzip sha256sum; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing+=("$cmd")
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        print_error "Missing required dependencies: ${missing[*]}"
        echo ""
        echo "Install instructions:"
        echo "  Ubuntu/Debian: sudo apt install curl jq unzip coreutils"
        echo "  macOS:         brew install curl jq unzip coreutils"
        echo "  Fedora/RHEL:   sudo dnf install curl jq unzip coreutils"
        exit 1
    fi
}

fetch_pypi_metadata() {
    print_info "Fetching PyPI metadata..."
    print_info "URL: $PYPI_JSON_URL"
    
    if ! JSON_DATA=$(curl -fsSL "$PYPI_JSON_URL" 2>/dev/null); then
        print_error "Failed to fetch package metadata from PyPI"
        echo ""
        echo "Possible causes:"
        echo "  - Network connectivity issues"
        echo "  - Invalid version: $EXTRACTOUS_VERSION"
        echo "  - PyPI service unavailable"
        exit 1
    fi
    
    # Validate JSON
    if ! echo "$JSON_DATA" | jq -e '.urls' >/dev/null 2>&1; then
        print_error "Invalid JSON response from PyPI"
        exit 1
    fi
    
    print_success "Metadata fetched successfully"
}

get_wheel_info() {
    local platform_tag="$1"
    
    echo "$JSON_DATA" | jq -r --arg tag "$platform_tag" '
        .urls[] | 
        select(.packagetype=="bdist_wheel") |
        select(.filename | contains($tag)) |
        "\(.filename)|\(.url)|\(.digests.sha256)"
    '
}

download_wheel() {
    local filename="$1"
    local url="$2"
    local sha256="$3"
    local wheel_path="$DOWNLOAD_DIR/$filename"
    
    # Check if already downloaded and verified
    if [ -f "$wheel_path" ]; then
        print_info "Wheel already downloaded, verifying checksum..."
        if echo "$sha256  $wheel_path" | sha256sum -c --status 2>/dev/null; then
            print_success "Using cached wheel"
            echo "$wheel_path"
            return 0
        else
            print_warning "Cached wheel checksum mismatch, re-downloading..."
            rm -f "$wheel_path"
        fi
    fi
    
    print_info "Downloading: $filename"
    
    # Download with progress bar if terminal supports it
    if [ -t 1 ]; then
        curl -fL --progress-bar -o "$wheel_path" "$url"
    else
        curl -fsSL -o "$wheel_path" "$url"
    fi
    
    # Verify checksum
    print_info "Verifying checksum..."
    if ! echo "$sha256  $wheel_path" | sha256sum -c --status; then
        print_error "Checksum verification failed!"
        rm -f "$wheel_path"
        exit 1
    fi
    
    print_success "Download complete and verified"
    echo "$wheel_path"
}

extract_libraries() {
    local wheel_path="$1"
    local platform_dir="$2"
    local extension="$3"
    
    local extract_subdir="$EXTRACT_DIR/$platform_dir"
    mkdir -p "$extract_subdir"
    
    print_info "Extracting wheel..."
    unzip -q "$wheel_path" -d "$extract_subdir"
    
    # Find and copy libraries
    local target_dir="$NATIVE_DIR/$platform_dir"
    mkdir -p "$target_dir"
    
    local copied=0
    local extractous_dir="$extract_subdir/extractous"
    
    if [ ! -d "$extractous_dir" ]; then
        print_warning "No extractous directory found in wheel"
        return 0
    fi
    
    # Copy all libraries with the correct extension
    while IFS= read -r -d '' lib; do
        local basename=$(basename "$lib")
        cp "$lib" "$target_dir/"
        print_success "  Copied: $basename ($(du -h "$lib" | cut -f1))"
        ((copied++))
    done < <(find "$extractous_dir" -type f -name "*.$extension*" -print0)
    
    echo "$copied"
}

process_platform() {
    local platform_tag="$1"
    local platform_dir="${PLATFORM_MAP[$platform_tag]}"
    local extension="${PLATFORM_EXTENSIONS[$platform_dir]}"
    
    print_section "Processing: $platform_dir (${platform_tag})"
    
    # Get wheel information
    local wheel_info=$(get_wheel_info "$platform_tag")
    
    if [ -z "$wheel_info" ]; then
        print_warning "No wheel found for platform: $platform_tag"
        return 1
    fi
    
    IFS='|' read -r filename url sha256 <<< "$wheel_info"
    
    print_info "Filename: $filename"
    print_info "Size: $(echo "$JSON_DATA" | jq -r --arg fn "$filename" '.urls[] | select(.filename==$fn) | .size' | numfmt --to=iec 2>/dev/null || echo "unknown")"
    
    # Download wheel
    local wheel_path=$(download_wheel "$filename" "$url" "$sha256")
    
    # Extract libraries
    local lib_count=$(extract_libraries "$wheel_path" "$platform_dir" "$extension")
    
    if [ "$lib_count" -gt 0 ]; then
        print_success "Extracted $lib_count libraries for $platform_dir"
        return 0
    else
        print_warning "No libraries found for $platform_dir"
        return 1
    fi
}

show_usage() {
    cat <<EOF
Usage: $0 [OPTIONS] [PLATFORMS...]

Extract native libraries from extractous Python wheels.

OPTIONS:
    -h, --help              Show this help message
    -v, --version VERSION   Specify extractous version (default: $EXTRACTOUS_VERSION)
    -a, --all               Extract for all supported platforms (default)
    -c, --clean             Clean existing native libraries before extraction

PLATFORMS:
    linux_amd64             Linux x86_64
    darwin_amd64            macOS Intel
    darwin_arm64            macOS Apple Silicon
    windows_amd64           Windows x86_64

EXAMPLES:
    # Extract for all platforms
    $0

    # Extract specific version
    $0 --version 0.2.1

    # Extract for specific platforms
    $0 linux_amd64 darwin_arm64

    # Clean and extract all
    $0 --clean --all

ENVIRONMENT:
    EXTRACTOUS_VERSION      Override default version
EOF
}

#=============================================================================
# Main Script
#=============================================================================

main() {
    local platforms_to_extract=()
    local clean_first=false
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -v|--version)
                EXTRACTOUS_VERSION="$2"
                PYPI_JSON_URL="https://pypi.org/pypi/${PACKAGE}/${EXTRACTOUS_VERSION}/json"
                shift 2
                ;;
            -a|--all)
                platforms_to_extract=("${!PLATFORM_MAP[@]}")
                shift
                ;;
            -c|--clean)
                clean_first=true
                shift
                ;;
            linux_amd64|darwin_amd64|darwin_arm64|windows_amd64)
                # Convert platform_dir to platform_tag
                for tag in "${!PLATFORM_MAP[@]}"; do
                    if [ "${PLATFORM_MAP[$tag]}" = "$1" ]; then
                        platforms_to_extract+=("$tag")
                        break
                    fi
                done
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                echo ""
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Default to all platforms if none specified
    if [ ${#platforms_to_extract[@]} -eq 0 ]; then
        platforms_to_extract=("${!PLATFORM_MAP[@]}")
    fi
    
    # Print header
    print_header "Extractous Native Library Extractor"
    echo "Version:      $EXTRACTOUS_VERSION"
    echo "Project Root: $PROJECT_ROOT"
    echo "Target Dir:   $NATIVE_DIR"
    echo "Platforms:    ${#platforms_to_extract[@]}"
    
    # Clean if requested
    if [ "$clean_first" = true ]; then
        print_section "Cleaning existing libraries"
        for platform_dir in "${PLATFORM_MAP[@]}"; do
            if [ -d "$NATIVE_DIR/$platform_dir" ]; then
                rm -rf "$NATIVE_DIR/$platform_dir"
                print_success "Cleaned: $platform_dir"
            fi
        done
    fi
    
    # Check dependencies
    print_section "Checking Dependencies"
    check_dependencies
    print_success "All dependencies available"
    
    # Fetch PyPI metadata
    print_section "Fetching Metadata"
    fetch_pypi_metadata
    
    # Process each platform
    local success_count=0
    local failed_platforms=()
    
    for platform_tag in "${platforms_to_extract[@]}"; do
        if process_platform "$platform_tag"; then
            ((success_count++))
        else
            failed_platforms+=("${PLATFORM_MAP[$platform_tag]}")
        fi
    done
    
    # Print summary
    print_header "Extraction Summary"
    
    echo ""
    echo "Results:"
    echo "  ✓ Successful:  $success_count / ${#platforms_to_extract[@]}"
    
    if [ ${#failed_platforms[@]} -gt 0 ]; then
        echo "  ✗ Failed:      ${#failed_platforms[@]}"
        echo ""
        echo "Failed platforms:"
        for platform in "${failed_platforms[@]}"; do
            echo "    - $platform"
        done
    fi
    
    echo ""
    echo "Native libraries structure:"
    if command -v tree >/dev/null 2>&1; then
        tree -L 2 "$NATIVE_DIR" || find "$NATIVE_DIR" -type f -name "*.*" | sed "s|$PROJECT_ROOT/||" | sort
    else
        find "$NATIVE_DIR" -type f \( -name "*.so" -o -name "*.dylib" -o -name "*.dll" \) | while read -r file; do
            size=$(du -h "$file" | cut -f1)
            echo "  $(echo "$file" | sed "s|$PROJECT_ROOT/||") ($size)"
        done
    fi
    
    echo ""
    if [ ${#failed_platforms[@]} -eq 0 ]; then
        print_success "All platforms extracted successfully!"
        echo ""
        echo "Next steps:"
        echo "  make build-ffi      # Build FFI for current platform"
        echo "  make build-ffi-all  # Build FFI for all platforms"
    else
        print_warning "Some platforms failed to extract"
        exit 1
    fi
}

# Run main function
main "$@"
