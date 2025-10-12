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

	// src/check_native_libs.go
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: Cannot determine file path. Please run the installer.")
		os.Exit(1)
	}

	projectRoot := filepath.Join(filepath.Dir(currentFile), "..")
	nativeDir := filepath.Join(projectRoot, "native")
	headerFile := filepath.Join(nativeDir, "include", "extractous.h")

	// Check if the header file exists.
	if _, err := os.Stat(headerFile); os.IsNotExist(err) {
		printError()
		os.Exit(1)
	}

	// Check for the platform-specific library directory.
	platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	libDir := filepath.Join(nativeDir, platform, "lib")

	if _, err := os.Stat(libDir); os.IsNotExist(err) {
		printError()
		os.Exit(1)
	}
}

func printError() {
	fmt.Fprintln(os.Stderr, "-------------------------------------------------------------------")
	fmt.Fprintln(os.Stderr, "Error: Native FFI libraries not found!")
	fmt.Fprintln(os.Stderr, "-------------------------------------------------------------------")
	fmt.Fprintln(os.Stderr, "This project uses CGO and requires pre-compiled native libraries")
	fmt.Fprintln(os.Stderr, "that were not found in the 'native/' directory.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "To fix this, please run the installer command from your project root:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "   go run github.com/rahulpoonia29/extractous-go/cmd/install@latest")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "This will download the correct libraries for your platform.")
	fmt.Fprintln(os.Stderr, "After running the installer, try your build again.")
	fmt.Fprintln(os.Stderr, "-------------------------------------------------------------------")
}
