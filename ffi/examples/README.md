# 📦 Extractous FFI Examples

This directory contains **C examples** demonstrating how to use the **Extractous FFI** library for text and metadata extraction from various document formats.

---

## Examples

| Example | Description |
|----------|-------------|
| **`basic.c`** | Simple file extraction with metadata |
| **`streaming.c`** | Stream large files without loading into memory |
| **`ocr.c`** | Extract scanned PDFs using OCR |

---

## Running Examples

```bash
# Basic extraction — extracts text and metadata from any supported document format
./basic document.pdf

# Streaming extraction — streams content from large files (>50MB)
./streaming large_document.pdf > output.txt

# OCR extraction — extracts text from scanned PDFs using Tesseract OCR
./ocr scanned_document.pdf
````

---

## Requirements

**Tesseract OCR** must be installed for OCR examples to work.

### Ubuntu / Debian

```bash
sudo apt install tesseract-ocr tesseract-ocr-eng
```

### macOS

```bash
brew install tesseract
```

### Windows

Download from the official repository: [https://github.com/UB-Mannheim/tesseract/wiki](https://github.com/UB-Mannheim/tesseract/wiki)

---

## Error Handling

All examples demonstrate robust error handling:

```c
int err = extractous_extractor_extract_file_to_string(...);
if (err != ERR_OK) {
    char* msg = extractous_error_message(err);
    fprintf(stderr, "Error: %s\n", msg);
    extractous_string_free(msg);
}
```

---

## Memory Management

Each example ensures proper cleanup of allocated resources:

```c
// Extract
extractous_extractor_extract_file_to_string(ex, path, &content, &meta);

// Use content and metadata
printf("%s\n", content);

// Cleanup
extractous_string_free(content);
extractous_metadata_free(meta);
extractous_extractor_free(ex);
```

---

## Common Issues

### Library Not Found

If you see:

```
error while loading shared libraries
```

Set the library path manually:

**Linux**

```bash
export LD_LIBRARY_PATH=../target/release:$LD_LIBRARY_PATH
./basic_extraction document.pdf
```

**macOS**

```bash
export DYLD_LIBRARY_PATH=../target/release:$DYLD_LIBRARY_PATH
./basic_extraction document.pdf
```

---

### OCR Not Available

If OCR examples fail with `ERR_OCR_NOT_AVAILABLE`:

1. **Install Tesseract:**

   ```bash
   # Ubuntu/Debian
   sudo apt install tesseract-ocr tesseract-ocr-eng

   # macOS
   brew install tesseract
   ```

2. **Verify installation:**

   ```bash
   tesseract --version
   ```

---

## Rough skeleton

```c
#include "../include/extractous.h"
#include <stdio.h>

int main(int argc, char** argv) {
    // 1. Create extractor
    CExtractor* ex = extractous_extractor_new();
    
    // 2. Configure (optional)
    ex = extractous_extractor_set_xml_output(ex, false);
    
    // 3. Extract
    char* content = NULL;
    CMetadata* meta = NULL;
    int err = extractous_extractor_extract_file_to_string(
        ex, "file.pdf", &content, &meta
    );
    
    // 4. Check error
    if (err != ERR_OK) {
        char* msg = extractous_error_message(err);
        fprintf(stderr, "Error: %s\n", msg);
        extractous_string_free(msg);
        extractous_extractor_free(ex);
        return 1;
    }
    
    // 5. Use results
    printf("%s\n", content);
    
    // 6. Cleanup
    extractous_string_free(content);
    extractous_metadata_free(meta);
    extractous_extractor_free(ex);
    
    return 0;
}
```
