<div align="center" style="margin-top: 20px">
  <h1>extractous-go</h1>
</div>

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/rahulpoonia29/extractous-go.svg)](https://pkg.go.dev/github.com/rahulpoonia29/extractous-go)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
<img src="https://img.shields.io/github/commit-activity/m/rahulpoonia29/extractous-go" alt="Commits per month">

</div>

<div align="center">

_Idiomatic Go bindings for [Extractous](https://github.com/yobix-ai/extractous), a high-performance library for extracting content and metadata from a wide variety of document formats._

</div>

---

This library brings the power and speed of the Rust-based `extractous` engine to the Go ecosystem. It uses CGO to wrap the native `extractous` FFI, providing a simple and efficient way to perform content extraction directly within your Go applications without relying on external services or APIs.

## Key Features

*   **High Performance**: Leverages the native speed and memory safety of the underlying Rust library.
*   **Simple Go API**: Provides an idiomatic Go interface for a seamless developer experience.
*   **Extensive Format Support**: Supports many file formats, including PDF, DOCX, XLSX, and more, by utilizing the power of Apache Tika compiled to a native library.
*   **OCR Capabilities**: Can extract text from images and scanned documents via Tesseract.
*   **Simple Installation**: Includes a command-line installer to automatically fetch the correct native libraries for your platform.

## Installation

Installation is a two-step process. First, you add the library to your project, and second, you run the installer to download the required native binaries.

**Step 1: Add the library to your project**

```bash
go get github.com/rahulpoonia29/extractous-go
```

**Step 2: Download the native libraries**

From the root of your project, run the installer. It will create a `native/` directory containing the platform-specific shared libraries.

```bash
go run github.com/rahulpoonia29/extractous-go/cmd/install@latest
```

The installer provides several flags for more advanced use cases, such as cross-compilation:

*   `--list-platforms`: List all available platforms from the latest release.
*   `--platform <name>`: Download libraries for a specific platform (e.g., `linux_amd64`).
*   `--all`: Download all available platforms.

## Quickstart

Here are a few examples of how to use `extractous-go`.

#### Extract content from a file to a string

```go
package main

import (
	"fmt"
	"log"

	"github.com/rahulpoonia29/extractous-go/src"
)

func main() {
	// Create a new extractor
	extractor := src.NewExtractor()

	// Extract text from a file
	content, metadata, err := extractor.ExtractFileToString("path/to/your/document.pdf")
	if err != nil {
		log.Fatalf("Failed to extract file: %v", err)
	}

	fmt.Println("--- Extracted Content ---")
	fmt.Println(content)

	fmt.Println("\n--- Metadata ---")
	for key, value := range metadata.AsMap() {
		fmt.Printf("%s: %s\n", key, value)
	}
}
```

#### Extract content to a stream for buffered reading

This is useful for large files to avoid loading the entire content into memory at once.

```go
package main

import (
	"fmt"
	"io"
	"log"

	"github.com/rahulpoonia29/extractous-go/src"
)

func main() {
	extractor := src.NewExtractor()

	// Get a reader for the extracted content
	reader, metadata, err := extractor.ExtractFile("path/to/your/document.docx")
	if err != nil {
		log.Fatalf("Failed to get reader: %v", err)
	}
	defer reader.Close()

	// Read the content in chunks
	buffer := make([]byte, 4096)
	content := ""
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read from stream: %v", err)
		}
		content += string(buffer[:n])
	}

	fmt.Println("--- Extracted Content (from stream) ---")
	fmt.Println(content)
	fmt.Printf("\n--- Metadata ---\n%v\n", metadata.AsMap())
}
```

## Distributing Your Application (e.g., with Wails)

Because this library uses CGO with **dynamically linked** native libraries (`.so`, `.dylib`, `.dll`), you must bundle these files with your final application. If you don't, your users will encounter runtime errors.

The process involves copying the libraries from the `native/` directory into your application package.

### Example: Wails v3 Build Configuration

If you are building a desktop app with Wails, you can automate this by modifying the platform-specific `Taskfile.yml` in your Wails project's `build/` directory.

**For macOS (`build/darwin/Taskfile.yml`)**

Add a task to copy the `.dylib` files into the `.app` bundle's `Contents/Frameworks` directory and make the `build` task depend on it.

```yaml
tasks:
  # ... other tasks
  copy-native-libs:
    summary: Copies the extractous-go native libraries.
    vars:
      LIB_PATH: '{{.ROOT_DIR}}/native/darwin_{{.ARCH}}/lib'
      TARGET_DIR: '{{.ROOT_DIR}}/build/bin/{{.APP_NAME}}.app/Contents/Frameworks'
    cmds:
      - mkdir -p {{.TARGET_DIR}}
      - cp {{.LIB_PATH}}/*.dylib {{.TARGET_DIR}}/
    preconditions:
      - sh: test -d {{.LIB_PATH}}
        msg: "Native libraries not found. Run 'go run github.com/rahulpoonia29/extractous-go/cmd/install@latest' first."

  build:
    deps:
      - task: common:build-frontend
      - task: copy-native-libs # Add this dependency
    # ... rest of build command
```

Similar modifications are needed for Linux (copying `.so` files and setting `rpath`) and Windows (copying `.dll` files next to the `.exe`). For detailed instructions, please refer to the guidance provided in the project's issue tracker or documentation.

## Roadmap (TODO)

This library is under active development. Future plans include:

- [ ] **Checksum Verification**: Add SHA256 checksum validation to the installer for improved security.
- [ ] **API Enhancement**: Expose more of the underlying `extractous` configuration options (e.g., for OCR, PDF parsing).
- [ ] **Static Linking**: Investigate providing a build option for static linking to simplify distribution.
- [ ] **More Examples**: Add more detailed examples, including for OCR and URL extraction.
- [ ] **Comprehensive Error Handling**: Improve error types to provide more context on failures.

## Acknowledgments

This project is a Go wrapper around the excellent [Extractous](https://github.com/yobix-ai/extractous) library. All credit for the high-performance extraction engine goes to its creators.

## 🕮 License

This project is licensed under the **Apache-2.0 license**. See the [LICENSE](LICENSE) file for details.
