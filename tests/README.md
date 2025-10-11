# Extractous-Go Test Suite

This directory contains comprehensive tests for the extractous-go library at multiple levels.

## Test Structure

```
tests/
├── ffi/                    # FFI layer tests (C interface)
│   ├── test_ffi_interface.c
│   └── Makefile
├── go/                     # Go binding tests
│   ├── bindings_test.go    # Unit tests for Go API
│   └── integration_test.go # Integration tests with actual files
└── testdata/               # Test files (created at runtime)
```

## Running Tests

### 1. FFI Layer Tests (C)

First, ensure the Rust FFI library is built:

```bash
cd ffi
cargo build --release
cd ..
```

Then run the FFI tests:

```bash
cd tests/ffi
make run
```

Or run individual test categories:

```bash
make run          # Run all tests
make clean        # Clean build artifacts
```

**What FFI tests validate:**
- Extractor lifecycle (new, free, double-free safety)
- Configuration functions (max length, encoding, XML output)
- PDF/Office/OCR configuration
- Error handling and null pointer safety
- URL extraction
- Memory management

### 2. Go Binding Tests

Run all Go tests:

```bash
cd tests/go
go test -v
```

Run specific test files:

```bash
go test -v -run TestExtractor    # Run extractor tests
go test -v -run TestPdfConfig    # Run PDF config tests
go test -v -run TestIntegration  # Run integration tests
```

Run with race detection:

```bash
go test -race -v
```

**What Go binding tests validate:**

#### `bindings_test.go` - Unit Tests
- Extractor lifecycle and nil-safety
- Configuration methods (max length, encoding, XML output)
- PDF/Office/OCR configuration
- Builder pattern and method chaining
- Error handling for nil extractors
- Metadata API (Get, GetAll, Has, Keys)
- CharSet constants

#### `integration_test.go` - Integration Tests
- Plain text file extraction
- Byte array extraction (string and stream)
- Configuration effects (max length, encoding, XML output)
- Metadata extraction and parsing
- Error handling (nonexistent files, empty files)
- Concurrent extraction (multiple goroutines)
- Multiple extractors on same file

## Test Data

Integration tests create temporary test files in `tests/testdata/` directory. These files are:
- Created at test runtime
- Cleaned up after each test
- Simple text files for validation

## Memory Management

**Important:** All config objects (PdfConfig, OfficeConfig, OcrConfig) use Go finalizers for automatic cleanup. You should **NOT** call any `Free()` method manually in Go code - they don't exist in the public API.

The FFI layer tests validate that the underlying C functions properly manage memory.

## Prerequisites

### For FFI Tests:
- GCC or compatible C compiler
- libextractous_ffi.so (built from Rust FFI layer)
- extractous.h header file

### For Go Tests:
- Go 1.25.1 or later
- CGo enabled
- libextractous_ffi.so in library path or proper LD_LIBRARY_PATH

## Troubleshooting

### FFI Tests

**Error: `libextractous_ffi.so: cannot open shared object file`**
```bash
# Ensure the library is built and in the right location
cd ffi && cargo build --release
# Check native/ directory for the compiled library
ls -la native/*/
```

**Error: `extractous.h: No such file or directory`**
```bash
# Regenerate the header with cbindgen
cd ffi
cbindgen --config cbindgen.toml --crate extractous-ffi --output ../include/extractous.h
```

### Go Tests

**Error: `undefined reference to extractous_*`**
- Ensure the FFI library is built: `cd ffi && cargo build --release`
- Check that CGo can find the library (see `src/cgo.go` for paths)

**Error: Package import issues**
- Ensure you're running tests from the `tests/go/` directory
- Module path should be `extractous-go` (check `go.mod`)

**Segmentation fault**
- This usually indicates a problem at the FFI boundary
- Run FFI tests first to validate the C interface
- Check that all CGo calls handle nil pointers correctly

## Continuous Integration

For CI pipelines, run tests in this order:

```bash
# 1. Build FFI library
cd ffi && cargo build --release && cd ..

# 2. Run FFI tests
cd tests/ffi && make run && cd ../..

# 3. Run Go tests
cd tests/go && go test -v -race && cd ../..
```

## Test Coverage

To generate coverage reports for Go tests:

```bash
cd tests/go
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

When adding new features:

1. **Add FFI tests first** - Validate the C interface in `tests/ffi/test_ffi_interface.c`
2. **Add Go unit tests** - Test the Go wrapper in `tests/go/bindings_test.go`
3. **Add integration tests** - Test end-to-end functionality in `tests/go/integration_test.go`

This ensures full validation from the C boundary up through the Go API.
