<div align="center" style="margin-top: 20px">
  <h1>Extractous Go</h1>
</div>

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/rahulpoonia29/extractous-go.svg)](https://pkg.go.dev/github.com/rahulpoonia29/extractous-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://github.com/rahulpoonia29/extractous-go/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/rahulpoonia29/extractous-go/actions/workflows/build.yml)

</div>

Go bindings for [Extractous](https://github.com/yobix-ai/extractous) — a fast, Rust-powered document extraction engine built on Apache Tika and Tesseract OCR.

---

## Features

- High-performance content extraction powered by Rust
- Support for 60+ file formats (PDF, DOCX, XLSX, PPTX, HTML, etc.)
- OCR capabilities for scanned documents and images via Tesseract
- Streaming API for large file processing with low memory usage
- Cross-platform support: Linux, macOS, Windows

---

## Installation

```bash
# Add the library
go get github.com/rahulpoonia29/extractous-go

# Download platform-specific native libraries
go run github.com/rahulpoonia29/extractous-go/cmd/install@latest
```

This installs the required native libraries under `native/`.
Use `--platform` or `--all` for cross-platform builds.

---

## Quick Start

### Extract Text from a File

```go
package main

import (
    "fmt"
    "log"
    "github.com/rahulpoonia29/extractous-go"
)

func main() {
    extractor := extractous.New()
    if extractor == nil {
		log.Fatal("Failed to create extractor")
	}
	defer extractor.Close()

    content, metadata, err := extractor.ExtractFileToString("document.pdf")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(content)
    fmt.Printf("Author: %s\n", metadata.Get("author"))
}
```

### Stream Large Files

```go
extractor := extractous.New()
if extractor == nil {
	log.Fatal("Failed to create extractor")
}
defer extractor.Close()

reader, metadata, err := extractor.ExtractFile("large_document.pdf")
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

buffer := make([]byte, 8192)
for {
    n, err := reader.Read(buffer)
    if err == io.EOF {
        break
    }
    if err != nil && err != io.EOF {
        log.Fatal(err)
    }
    // Process buffer[:n]
}
```

See ([examples](./examples)) for more usage patterns.

### Advanced Configuration

```go
pdfConfig := extractous.NewPdfConfig().
    SetOcrStrategy(extractous.PdfOcrAuto).
    SetExtractInlineImages(true)

extractor := extractous.New().
    SetPdfConfig(pdfConfig)

content, _, err := extractor.ExtractFileToString("document.pdf")
```

---

## Distribution

This library uses CGO and dynamically linked native libraries.
When distributing applications, bundle the required shared libraries (`.so`, `.dylib`, `.dll`) with your executable.

---

## Building & Running

### 1. Linux

```bash
# Enable CGO and set compiler
export CGO_ENABLED=1
export CC=gcc
export CXX=g++

# Set header and library paths
export CGO_CFLAGS="-I$(pwd)/native/linux_amd64/include"
export CGO_LDFLAGS="-L$(pwd)/native/linux_amd64/lib -lextractous_ffi -lstdc++ -Wl,-rpath,$(pwd)/native/linux_amd64/lib"

# Build and run
go build main.go
./main
```

### 2. macOS

```bash
# Enable CGO and set compiler
export CGO_ENABLED=1
export CC=clang
export CXX=clang++

# Set header and library paths
export CGO_CFLAGS="-I$(pwd)/native/darwin_arm64/include"
export CGO_LDFLAGS="-L$(pwd)/native/darwin_arm64/lib -lextractous_ffi -lc++ -Wl,-rpath,$(pwd)/native/darwin_arm64/lib"

# Build and run
go build main.go
./main
```

### 3. Windows (PowerShell)

```powershell
# Enable CGO and set compiler
$env:CGO_ENABLED = 1
$env:CC = "C:\msys64\mingw64\bin\gcc.exe"
$env:CXX = "C:\msys64\mingw64\bin\g++.exe"
$env:Path = "C:\msys64\mingw64\bin;" + $env:Path

# Set header and library paths
$env:CGO_CFLAGS = "-I$(pwd)\native\windows_amd64\include"
$env:CGO_LDFLAGS = "-L$(pwd)\native\windows_amd64\lib -lextractous_ffi"

# Build the executable
go build main.go

# Copy the binary next to the DLLs and run
copy main.exe native\windows_amd64\lib\
cd native\windows_amd64\lib
.\main.exe
```

---

## Performance

| Library           | Throughput (MB/s) | Memory (MB) | Accuracy (%) |
| ----------------- | ----------------- | ----------- | ------------ |
| extractous-string | 36.70             | 15.78       | 86.95        |
| extractous-stream | 14.16             | 21.83       | 87.74        |
| ledongthuc-pdf    | 79.38             | 44.67       | 82.02        |

---

## Supported Formats

- PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX
- HTML, XML, Markdown, TXT, RTF, CSV
- OpenDocument (ODT, ODS, ODP)
- And many more...

See [Apache Tika Supported Formats](https://tika.apache.org/2.0.0/formats.html) for the full list.

---

## Requirements

- Go 1.19 or later
- CGO enabled
- Platform-specific native libraries (installed via the `install` command)
- Tesseract (only required for OCR functionality)

---

## Roadmap

- [ ] Complete Windows custom loader implementation
- [ ] Add checksum verification to installer
- [ ] Add tests for FFI layer and Go bindings

---

## License

Licensed under the Apache License 2.0 — see [LICENSE](LICENSE) for details.
