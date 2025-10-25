# Extractous Go

Go bindings for [Extractous](https://github.com/yobix-ai/extractous) - fast, high-performance, rust-powered document extraction built on Apache Tika and Tesseract OCR.

[![Go Reference](https://pkg.go.dev/badge/github.com/rahulpoonia29/extractous-go.svg)](https://pkg.go.dev/github.com/rahulpoonia29/extractous-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://github.com/rahulpoonia29/extractous-go/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/rahulpoonia29/extractous-go/actions/workflows/build.yml)

---

## Features

- **High Performance**: Built on Rust for maximum throughput and minimal memory overhead
- **60+ File Formats**: PDF, Office documents (DOCX, XLSX, PPTX), HTML, XML, and more
- **OCR Support**: Extract text from scanned documents and images using Tesseract
- **Streaming API**: Process large files with minimal memory usage

---

## Installation

### Step 1: Install the Go Package

```bash
go get github.com/rahulpoonia29/extractous-go
```

### Step 2: Download Native Libraries

```bash
# Download libraries for your current platform
go run github.com/rahulpoonia29/extractous-go/cmd/install@latest

# Download for a specific platform
go run github.com/rahulpoonia29/extractous-go/cmd/install@latest --platform linux-amd64

# Download for all platforms (useful for cross-compilation)
go run github.com/rahulpoonia29/extractous-go/cmd/install@latest --all
```

This creates a `native/` directory with libraries for your specific platform.

---

## Quick Start

### Basic Text Extraction

```go
package main

import (
    "fmt"
    "log"

    "github.com/rahulpoonia29/extractous-go"
)

func main() {
    // Create a new extractor
    extractor := extractous.New()
    if extractor == nil {
        log.Fatal("Failed to create extractor")
    }
    defer extractor.Close()

    // Extract text and metadata from file
    content, metadata, err := extractor.ExtractFileToString("document.pdf")
    if err != nil {
        log.Fatalf("Extraction failed: %v", err)
    }

    // Results
    fmt.Println("Content:", content)
    fmt.Println("Metadata:", metadata)
}
```

### Streaming Large Files

For memory efficient processing of large documents:

```go
package main

import (
    "fmt"
    "io"
    "log"

    "github.com/rahulpoonia29/extractous-go"
)

func main() {
    extractor := extractous.New()
    if extractor == nil {
        log.Fatal("Failed to create extractor")
    }
    defer extractor.Close()

    // Get a streaming reader for the document
    reader, metadata, err := extractor.ExtractFile("large_document.pdf")
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    // Process the document in chunks
    buffer := make([]byte, 8192)
    for {
        n, err := reader.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        // Process buffer[:n]
        fmt.Printf("Read %d bytes\n", n)
    }
}
```

### Configuration

```go
package main

import (
    "log"

    "github.com/rahulpoonia29/extractous-go"
)

func main() {
    // Configure PDF extraction with OCR
    pdfConfig := extractous.NewPdfConfig()
    pdfConfig.SetOcrStrategy(extractous.PdfOcrAuto)
    pdfConfig.SetExtractInlineImages(true)
    pdfConfig.SetExtractAnnotationText(true)

    // Configure OCR settings
    ocrConfig := extractous.NewTesseractOcrConfig()
    ocrConfig.SetLanguage("eng")
    ocrConfig.SetDensity(300)

    // Apply configurations to extractor
    extractor := extractous.New()
    extractor.SetPdfConfig(pdfConfig)
    extractor.SetTesseractOcrConfig(ocrConfig)
    extractor.SetXmlOutput(true) // Enable structured output

    defer extractor.Close()

    content, _, err := extractor.ExtractFileToString("scanned_document.pdf")
    if err != nil {
        log.Fatal(err)
    }

    log.Println(content)
}
```

---

## Building

The library uses CGO to interface with native libraries. Below are platform-specific build instructions.

### Prerequisite

- CGO enabled
- Native libraries
- Platform-specific C compiler

### Linux and macOS

```bash
# Set up environment
export CGO_ENABLED=1
export CC=gcc
export CXX=g++

# Set library path for the build
export CGO_LDFLAGS="-L$(pwd)/native/$(go env GOOS)_$(go env GOARCH) -lextractous_ffi"

# Build the application
go build -o myapp main.go

# Set the library path for runtime before executing
export LD_LIBRARY_PATH="$(pwd)/native/$(go env GOOS)_$(go env GOARCH):$LD_LIBRARY_PATH" # For Linux
export DYLD_LIBRARY_PATH="$(pwd)/native/$(go env GOOS)_$(go env GOARCH):$DYLD_LIBRARY_PATH" # For macOS

./myapp

```

### Windows (PowerShell)

```powershell
# Set up environment
$env:CGO_ENABLED = "1"
$env:CC = "gcc"
$env:CXX = "g++"

# Set library path for the build
$env:CGO_LDFLAGS = "-L$pwd\native\windows_amd64 -lextractous_ffi" # Only x86-64 is supported

# Build the application
go build -o myapp.exe main.go

# Add the DLL to the system's path
$env:Path = "$pwd\native\windows_amd64;" + $env:Path
.\myapp.exe
```

---

## Error Handling

### Basic Error Handling

```go
content, metadata, err := extractor.ExtractFileToString("document.pdf")
if err != nil {
    // Check error type
    if errors.Is(err, extractous.ErrIO) {
        log.Println("File I/O error")
    } else if errors.Is(err, extractous.ErrExtraction) {
        log.Println("Document extraction failed")
    }

    log.Fatal(err)
}
```

### Error Handling with Debug Info

```go
content, metadata, err := extractor.ExtractFileToString("document.pdf")
if err != nil {
    // Get structured error information
    var extractErr *extractous.ExtractError
    if errors.As(err, &extractErr) {
        fmt.Printf("Error code: %d\n", extractErr.Code)
        fmt.Printf("Message: %s\n", extractErr.Message)

        // Optionally get detailed debug information
        // (includes full error chain and backtrace if available)
        if debug := extractErr.Debug(); debug != "" {
            fmt.Printf("Debug info:\n%s\n", debug)
        }
    }
}
```
---

## Performance

| Operation          | Throughput (MB/s) | Memory (MB) | Accuracy (%) |
| ------------------ | ----------------- | ----------- | ------------ |
| String Extraction  | 36.70             | 15.78       | 86.95        |
| Stream Extraction  | 14.16             | 21.83       | 87.74        |
| Reference (Go PDF) | 79.38             | 44.67       | 82.02        |

---

## Supported Formats

Extractous Go supports PDF, Microsoft Office, OpenDocument, HTML/XML, plain text, images (with OCR) and more.

For the full list of supported formats, see [Apache Tika Supported Formats](https://tika.apache.org/2.0.0/formats.html).

---

## Requirements

### Runtime Requirements

- Go 1.19 or later
- CGO enabled (`CGO_ENABLED=1`)
- Platform-specific native libraries (provided by installer)
- **Tesseract OCR**: Required only for OCR functionality on images and scanned PDFs
  - Ubuntu/Debian: `sudo apt-get install tesseract-ocr`
  - macOS: `brew install tesseract`
  - Windows: Download from [Tesseract at UB Mannheim](https://github.com/UB-Mannheim/tesseract/wiki)

---

## Distribution

When distributing applications built with extractous-go:

1. **Bundle Native Libraries**: Include the platform-specific `.so`, `.dylib`, or `.dll` files with your application.

2. **Set Library Search Path**:
   - **Linux**: Set `LD_LIBRARY_PATH` or install to `/usr/local/lib`
   - **macOS**: Set `DYLD_LIBRARY_PATH` or use `@rpath`
   - **Windows**: Place DLL in the same directory as the executable or in `System32`

3. **Cross-Platform Builds**: Download libraries for all target platforms using:
   ```bash
   go run github.com/rahulpoonia29/extractous-go/cmd/install@latest --all
   ```

---

## Acknowledgments

- [Extractous](https://github.com/yobix-ai/extractous) - The underlying Rust library
- [Apache Tika](https://tika.apache.org/) - Document extraction engine
- [Tesseract OCR](https://github.com/tesseract-ocr/tesseract) - OCR engine
