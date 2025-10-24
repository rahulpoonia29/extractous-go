//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	// This script is run by `go generate` to check if the native libraries
	// required for CGO exist. If they don't, it prints a helpful error
	// message to guide the user.

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: Cannot determine file path. Please run the installer.")
		os.Exit(1)
	}

	projectRoot := filepath.Dir(currentFile)
	nativeDir := filepath.Join(projectRoot, "native")

	// Check if native directory exists
	if _, err := os.Stat(nativeDir); os.IsNotExist(err) {
		printError()
		os.Exit(1)
	}

	// Check for the platform-specific library directory
	platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	libDir := filepath.Join(nativeDir, platform)

	if _, err := os.Stat(libDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Error: Native libraries not found for %s!\n", platform)
		fmt.Fprintf(os.Stderr, "Expected library directory: %s\n", libDir)
		fmt.Fprintf(os.Stderr, "\n")
		printError()
		os.Exit(1)
	}

	// Verify the actual library file exists
	var libFile string
	switch runtime.GOOS {
	case "windows":
		libFile = filepath.Join(libDir, "extractous_ffi.dll")
	case "darwin":
		libFile = filepath.Join(libDir, "libextractous_ffi.dylib")
	default: // linux
		libFile = filepath.Join(libDir, "libextractous_ffi.so")
	}

	if _, err := os.Stat(libFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Library file not found: %s\n", libFile)
		printError()
		os.Exit(1)
	}

	// Success - libraries are present
	fmt.Printf("✓ Native libraries verified for %s\n", platform)
	fmt.Printf("  Library: %s\n", libFile)
}

func printError() {
	fmt.Fprintln(os.Stderr, "Error: Native FFI libraries not found!")
	fmt.Fprintln(os.Stderr, "This project uses CGO and requires pre-compiled native libraries")
	fmt.Fprintln(os.Stderr, "that were not found in the 'native/' directory.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "To fix this, please run the installer command from your project root:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  go run github.com/rahulpoonia29/extractous-go/cmd/install@latest")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "This will download the correct libraries for your platform.")
	fmt.Fprintln(os.Stderr, "After running the installer, try your build again.")
}
