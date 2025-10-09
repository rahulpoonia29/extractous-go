#!/usr/bin/env bash
set -euo pipefail

#=============================================================================
# Extractous-Go Central Build Script
#=============================================================================
# One script to rule them all - orchestrates the entire build process
#=============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m'

# Configuration
EXTRACTOUS_VERSION="${EXTRACTOUS_VERSION:-0.3.0}"
CARGO_PROFILE="${CARGO_PROFILE:-release}"
DEFAULT_PLATFORMS=("linux_amd64" "darwin_amd64" "darwin_arm64" "windows_amd64")

#=============================================================================
# Helper Functions
#=============================================================================

print_banner() {
    echo -e "${CYAN}${BOLD}"
    cat <<'EOF'
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║              EXTRACTOUS-GO BUILD SYSTEM                       ║
║                                                               ║
║   Cross-Platform Document Extraction FFI Builder              ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
EOF
    echo -e "${NC}"
}

print_header() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}${BOLD}  $1${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
}

print_section() {
    echo ""
    echo -e "${BLUE}───────────────────────────────────────────────────────────────${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}───────────────────────────────────────────────────────────────${NC}"
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

print_step() {
    echo -e "${MAGENTA}▶${NC} ${BOLD}$1${NC}"
}

show_usage() {
    cat <<EOF
${BOLD}USAGE:${NC}
    $0 [COMMAND] [OPTIONS]

${BOLD}COMMANDS:${NC}
    ${BOLD}setup${NC}               Complete initial setup (toolchains + wheels)
    ${BOLD}extract${NC}             Extract native libraries from Python wheels
    ${BOLD}build${NC}               Build FFI library (current platform)
    ${BOLD}build-all${NC}           Build FFI for all platforms
    ${BOLD}test${NC}                Run C smoke tests
    ${BOLD}clean${NC}               Clean build artifacts
    ${BOLD}clean-all${NC}           Clean everything including native libraries
    ${BOLD}check${NC}               Verify environment setup
    ${BOLD}status${NC}              Show build status
    ${BOLD}help${NC}                Show this help message

${BOLD}OPTIONS:${NC}
    -v, --version VERSION      Extractous version (default: $EXTRACTOUS_VERSION)
    -p, --platform PLATFORM    Target platform (for build command)
    -P, --profile PROFILE      Cargo build profile (release/debug, default: release)
    --clean                    Clean before building
    --verbose                  Verbose output
    -h, --help                 Show this help

${BOLD}PLATFORMS:${NC}
    linux_amd64               Linux x86_64
    darwin_amd64              macOS Intel
    darwin_arm64              macOS Apple Silicon
    windows_amd64             Windows x86_64

${BOLD}EXAMPLES:${NC}
    # Complete setup from scratch
    $0 setup

    # Build for current platform
    $0 build

    # Build for all platforms
    $0 build-all

    # Build for specific platform
    $0 build --platform darwin_arm64

    # Extract specific version
    $0 extract --version 0.2.1

    # Clean build
    $0 build --clean

    # Check environment
    $0 check

    # Show build status
    $0 status

${BOLD}ENVIRONMENT VARIABLES:${NC}
    EXTRACTOUS_VERSION        Override default version
    CARGO_PROFILE             Override build profile (release/debug)

${BOLD}TYPICAL WORKFLOW:${NC}
    1. $0 setup              # One-time setup
    2. $0 build              # Build for development
    3. $0 test               # Run tests
    4. $0 build-all          # Build for release
EOF
}

run_command() {
    local cmd="$1"
    shift
    print_step "Running: $cmd"
    if "$cmd" "$@"; then
        print_success "Command completed: $cmd"
        return 0
    else
        print_error "Command failed: $cmd"
        return 1
    fi
}

#=============================================================================
# Command Implementations
#=============================================================================

cmd_setup() {
    print_header "Complete Setup"
    
    print_section "Step 1: Installing Rust Toolchains"
    run_command "$SCRIPT_DIR/setup_toolchains.sh"
    
    print_section "Step 2: Extracting Native Libraries"
    run_command "$SCRIPT_DIR/extract_wheels.sh"
    
    print_section "Step 3: Verifying Environment"
    run_command "$SCRIPT_DIR/check_env.sh"
    
    print_header "Setup Complete!"
    echo ""
    print_success "Environment is ready for building"
    echo ""
    echo "Next steps:"
    echo "  $0 build        # Build for current platform"
    echo "  $0 build-all    # Build for all platforms"
    echo "  $0 test         # Run tests"
}

cmd_extract() {
    print_header "Extracting Native Libraries"
    
    local extract_args=()
    
    if [ -n "${OPT_VERSION:-}" ]; then
        extract_args+=("--version" "$OPT_VERSION")
    fi
    
    if [ -n "${OPT_PLATFORM:-}" ]; then
        extract_args+=("$OPT_PLATFORM")
    fi
    
    if [ "${OPT_CLEAN:-false}" = true ]; then
        extract_args+=("--clean")
    fi
    
    run_command "$SCRIPT_DIR/extract_wheels.sh" "${extract_args[@]}"
}

cmd_build() {
    print_header "Building FFI Library"
    
    local platform="${OPT_PLATFORM:-}"
    
    if [ "${OPT_CLEAN:-false}" = true ]; then
        print_section "Cleaning Build Artifacts"
        cd "$PROJECT_ROOT/ffi"
        cargo clean
        print_success "Clean complete"
    fi
    
    if [ -n "$platform" ]; then
        print_info "Target Platform: $platform"
        run_command "$SCRIPT_DIR/build_ffi.sh" "$platform"
    else
        local current_platform=$("$SCRIPT_DIR/detect_platform.sh")
        print_info "Target Platform: $current_platform (detected)"
        run_command "$SCRIPT_DIR/build_ffi.sh" "$current_platform"
    fi
}

cmd_build_all() {
    print_header "Building FFI for All Platforms"
    
    if [ "${OPT_CLEAN:-false}" = true ]; then
        print_section "Cleaning Build Artifacts"
        cd "$PROJECT_ROOT/ffi"
        cargo clean
        print_success "Clean complete"
    fi
    
    run_command "$SCRIPT_DIR/build_ffi_all.sh"
}

cmd_test() {
    print_header "Running Tests"
    
    print_section "Building Current Platform"
    local current_platform=$("$SCRIPT_DIR/detect_platform.sh")
    run_command "$SCRIPT_DIR/build_ffi.sh" "$current_platform"
    
    print_section "Running C Smoke Tests"
    run_command "$SCRIPT_DIR/test_ffi.sh"
}

cmd_clean() {
    print_header "Cleaning Build Artifacts"
    
    print_section "Cleaning Rust Build"
    cd "$PROJECT_ROOT/ffi"
    cargo clean
    print_success "Rust artifacts cleaned"
    
    print_section "Cleaning Generated Header"
    if [ -f "$PROJECT_ROOT/include/extractous.h" ]; then
        rm -f "$PROJECT_ROOT/include/extractous.h"
        print_success "Header cleaned"
    fi
    
    print_section "Cleaning Test Artifacts"
    if [ -d "$PROJECT_ROOT/tests/smoke" ]; then
        make -C "$PROJECT_ROOT/tests/smoke" clean 2>/dev/null || true
        print_success "Test artifacts cleaned"
    fi
    
    print_success "Clean complete"
}

cmd_clean_all() {
    print_header "Deep Clean"
    
    # Run normal clean
    cmd_clean
    
    print_section "Removing Native Libraries"
    local removed=0
    for platform in "${DEFAULT_PLATFORMS[@]}"; do
        local native_dir="$PROJECT_ROOT/native/$platform"
        if [ -d "$native_dir" ]; then
            rm -rf "$native_dir"
            print_success "Cleaned: $platform"
            ((removed++))
        fi
    done
    
    if [ $removed -eq 0 ]; then
        print_info "No native libraries to clean"
    fi
    
    print_success "Deep clean complete"
    echo ""
    echo "To rebuild from scratch, run:"
    echo "  $0 setup"
}

cmd_check() {
    print_header "Environment Check"
    run_command "$SCRIPT_DIR/check_env.sh"
}

cmd_status() {
    print_header "Build Status"
    
    # Check Rust installation
    print_section "Rust Environment"
    if command -v rustc &> /dev/null; then
        print_success "Rust: $(rustc --version)"
        print_success "Cargo: $(cargo --version)"
    else
        print_error "Rust: Not installed"
    fi
    
    # Check native libraries
    print_section "Native Libraries"
    local lib_count=0
    for platform in "${DEFAULT_PLATFORMS[@]}"; do
        local native_dir="$PROJECT_ROOT/native/$platform"
        if [ -d "$native_dir" ]; then
            local count=$(find "$native_dir" -type f \( -name "*.so" -o -name "*.dylib" -o -name "*.dll" \) 2>/dev/null | wc -l)
            if [ "$count" -gt 0 ]; then
                print_success "$platform: $count libraries"
                ((lib_count++))
            else
                print_warning "$platform: Directory exists but empty"
            fi
        else
            print_warning "$platform: Not extracted"
        fi
    done
    
    if [ $lib_count -eq 0 ]; then
        echo ""
        print_warning "No native libraries found. Run: $0 extract"
    fi
    
    # Check FFI build
    print_section "FFI Libraries"
    local ffi_count=0
    for platform in "${DEFAULT_PLATFORMS[@]}"; do
        local native_dir="$PROJECT_ROOT/native/$platform"
        local lib_name=""
        
        case "$platform" in
            linux_amd64)   lib_name="libextractous_ffi.so" ;;
            darwin_*)      lib_name="libextractous_ffi.dylib" ;;
            windows_amd64) lib_name="extractous_ffi.dll" ;;
        esac
        
        if [ -f "$native_dir/$lib_name" ]; then
            local size=$(du -h "$native_dir/$lib_name" | cut -f1)
            print_success "$platform: $lib_name ($size)"
            ((ffi_count++))
        else
            print_warning "$platform: FFI not built"
        fi
    done
    
    if [ $ffi_count -eq 0 ]; then
        echo ""
        print_warning "No FFI libraries built. Run: $0 build"
    fi
    
    # Check generated header
    print_section "Generated Files"
    if [ -f "$PROJECT_ROOT/include/extractous.h" ]; then
        print_success "C Header: include/extractous.h"
    else
        print_warning "C Header: Not generated"
    fi
    
    # Summary
    print_section "Summary"
    echo ""
    if [ $lib_count -eq ${#DEFAULT_PLATFORMS[@]} ] && [ $ffi_count -eq ${#DEFAULT_PLATFORMS[@]} ]; then
        print_success "All platforms ready!"
    elif [ $lib_count -eq ${#DEFAULT_PLATFORMS[@]} ]; then
        print_warning "Native libraries ready, but FFI needs building"
        echo ""
        echo "Run: $0 build-all"
    elif [ $lib_count -eq 0 ]; then
        print_warning "No native libraries extracted"
        echo ""
        echo "Run: $0 setup"
    else
        print_warning "Partial build - some platforms missing"
        echo ""
        echo "Run: $0 setup"
    fi
}

#=============================================================================
# Main Script
#=============================================================================

main() {
    # Parse options
    local command=""
    local OPT_VERSION=""
    local OPT_PLATFORM=""
    local OPT_PROFILE="$CARGO_PROFILE"
    local OPT_CLEAN=false
    local OPT_VERBOSE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            setup|extract|build|build-all|test|clean|clean-all|check|status|help)
                command="$1"
                shift
                ;;
            -v|--version)
                OPT_VERSION="$2"
                EXTRACTOUS_VERSION="$2"
                export EXTRACTOUS_VERSION
                shift 2
                ;;
            -p|--platform)
                OPT_PLATFORM="$2"
                shift 2
                ;;
            -P|--profile)
                OPT_PROFILE="$2"
                CARGO_PROFILE="$2"
                export CARGO_PROFILE
                shift 2
                ;;
            --clean)
                OPT_CLEAN=true
                shift
                ;;
            --verbose)
                OPT_VERBOSE=true
                set -x
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo ""
                show_usage
                exit 1
                ;;
        esac
    done
    
    # Show banner
    print_banner
    
    # Default command
    if [ -z "$command" ]; then
        show_usage
        exit 0
    fi
    
    # Change to project root
    cd "$PROJECT_ROOT"
    
    # Execute command
    case "$command" in
        setup)
            cmd_setup
            ;;
        extract)
            cmd_extract
            ;;
        build)
            cmd_build
            ;;
        build-all)
            cmd_build_all
            ;;
        test)
            cmd_test
            ;;
        clean)
            cmd_clean
            ;;
        clean-all)
            cmd_clean_all
            ;;
        check)
            cmd_check
            ;;
        status)
            cmd_status
            ;;
        help)
            show_usage
            ;;
        *)
            print_error "Unknown command: $command"
            show_usage
            exit 1
            ;;
    esac
}

# Run main
main "$@"
