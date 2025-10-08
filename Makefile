.PHONY: all build clean test install extract-libs copy-libs

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
    OS = linux
endif
ifeq ($(UNAME_S),Darwin)
    OS = darwin
endif
ifeq ($(UNAME_M),x86_64)
    ARCH = amd64
endif
ifeq ($(UNAME_M),arm64)
    ARCH = arm64
endif
ifeq ($(UNAME_M),aarch64)
    ARCH = arm64
endif

PLATFORM = $(OS)_$(ARCH)

all: build

# Extract native libs from Python wheels
extract-libs:
	@echo "Extracting native libraries from Python wheels..."
	./scripts/extract_wheels.sh

# Build Rust FFI library for current platform
build-ffi:
	@echo "Building Rust FFI for $(PLATFORM)..."
	cd extractous-ffi && cargo build --release

# Build for all platforms (requires cross-compilation setup)
build-all:
	@echo "Building for all platforms..."
	./scripts/build_all.sh

# Copy built libraries to Go native/ directory
copy-libs:
	@echo "Copying libraries to native/$(PLATFORM)..."
	./scripts/copy_libs.sh $(PLATFORM)

# Build Go package
build: build-ffi copy-libs
	@echo "Building Go package..."
	go build -v ./...

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	cd extractous-ffi && cargo clean
	rm -rf native/
	go clean

# Install Go package
install:
	go install ./...

# Generate Go documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060
