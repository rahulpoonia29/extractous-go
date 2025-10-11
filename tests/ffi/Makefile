# Makefile for FFI Tests

CC = gcc
CFLAGS = -Wall -Wextra -I../../include
LDFLAGS = \
    -L../../extractous-ffi/target/release \
    -L../../native/$(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') \
    -Wl,-rpath,../../extractous-ffi/target/release \
    -Wl,-rpath,../../native/$(shell uname -s | tr '[:upper:]' '[:lower:]')_$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
LIBS = -lextractous_ffi -ldl -lm -lpthread

TEST_BINS = test_ffi_interface

.PHONY: all clean run

all: $(TEST_BINS)

test_ffi_interface: test_ffi_interface.c
	$(CC) $(CFLAGS) -o $@ $< $(LDFLAGS) $(LIBS)

run: all
	@echo "Running FFI interface tests..."
	@./test_ffi_interface

clean:
	rm -f $(TEST_BINS)
	rm -f *.o

help:
	@echo "FFI Test Suite"
	@echo ""
	@echo "Targets:"
	@echo "  all     - Build all test binaries"
	@echo "  run     - Build and run all tests"
	@echo "  clean   - Remove test binaries"
	@echo "  help    - Show this help message"
